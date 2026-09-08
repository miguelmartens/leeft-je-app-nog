// Demo-app voor de workshop "Leeft je app nog? - een introductie tot probes
// in Kubernetes".
//
// Bewust één bestand en zonder externe dependencies, zodat het op een projector
// leesbaar blijft. Met de /debug-endpoints maak je de app tijdens de sessie
// stuk en weer heel, zonder opnieuw te bouwen of uit te rollen.
//
//	GET  /startupz               ben ik opgestart?
//	GET  /healthz                leef ik nog?          (startup + liveness probe)
//	GET  /readyz                 mag ik verkeer?       (readiness probe)
//	GET  /work                   het echte werk        (het "gebruikersverkeer")
//	POST /debug/unready          readiness aan/uit
//	POST /debug/deadlock         /healthz laten hangen, liveness gaat falen
//	POST /debug/dependency-down  de nagebootste database aan/uit
//
// Instelbaar via de omgeving: PORT (8080), STARTUP_DELAY (5s), WORK_DURATION (100ms).
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Na SIGTERM krijgen lopende verzoeken zoveel tijd om af te ronden. Houd dit
// onder terminationGracePeriodSeconds min de duur van de preStop-hook, anders
// maakt SIGKILL er alsnog een eind aan.
const afsluitTimeout = 25 * time.Second

// app is de complete toestand van de demo. Alles is atomic, want de
// HTTP-handlers lezen en schrijven vanuit meerdere goroutines tegelijk.
type app struct {
	opgestart      atomic.Bool // is het opstarten klaar?
	afsluiten      atomic.Bool // is SIGTERM binnengekomen?
	unready        atomic.Bool // readiness handmatig uitgezet
	deadlock       atomic.Bool // /healthz blijft hangen
	dependencyDown atomic.Bool // de nagebootste database is onbereikbaar

	pod      string        // podnaam, zodat je in de output ziet wie antwoordt
	werkduur time.Duration // hoelang /work doet alsof het iets nuttigs doet
}

// startupz beantwoordt de vraag "ben ik al opgestart?". Het manifest laat de
// startup probe naar /healthz wijzen, net als in de presentatie; dit endpoint
// geeft hetzelfde antwoord en is handig om met de hand te controleren.
func (a *app) startupz(w http.ResponseWriter, r *http.Request) {
	if !a.opgestart.Load() {
		a.antwoord(w, http.StatusServiceUnavailable, "bezig met opstarten")
		return
	}
	a.antwoord(w, http.StatusOK, "opgestart")
}

// healthz beantwoordt de liveness probe: leef ik nog? Deze kijkt alleen naar
// zichzelf en nooit naar een database of een andere service. Zodra dat wel
// gebeurt, herstart één trage dependency al je replica's tegelijk.
func (a *app) healthz(w http.ResponseWriter, r *http.Request) {
	if a.deadlock.Load() {
		// Nagebootste deadlock: we antwoorden simpelweg niet meer. De kubelet
		// loopt tegen timeoutSeconds aan en herstart na failureThreshold keer
		// de container.
		<-r.Context().Done()
		return
	}
	if !a.opgestart.Load() {
		a.antwoord(w, http.StatusServiceUnavailable, "bezig met opstarten")
		return
	}
	a.antwoord(w, http.StatusOK, "leeft")
}

// readyz beantwoordt de readiness probe: mag ik verkeer? Deze mag wel naar
// buiten kijken. Faalt hij, dan gaat dit pod uit de EndpointSlice en blijft de
// container gewoon draaien.
func (a *app) readyz(w http.ResponseWriter, r *http.Request) {
	switch {
	case !a.opgestart.Load():
		a.antwoord(w, http.StatusServiceUnavailable, "bezig met opstarten")
	case a.afsluiten.Load():
		a.antwoord(w, http.StatusServiceUnavailable, "aan het afsluiten")
	case a.unready.Load():
		a.antwoord(w, http.StatusServiceUnavailable, "handmatig op unready gezet")
	case a.dependencyDown.Load():
		a.antwoord(w, http.StatusServiceUnavailable, "dependency onbereikbaar")
	default:
		a.antwoord(w, http.StatusOK, "klaar voor verkeer")
	}
}

// work is het gebruikersverkeer: het endpoint waar de load-generator op schiet.
// Tijdens het opstarten is er nog niets te halen - dat is precies wat je zonder
// readiness probe aan je gebruikers laat zien.
func (a *app) work(w http.ResponseWriter, r *http.Request) {
	if !a.opgestart.Load() {
		a.antwoord(w, http.StatusServiceUnavailable, "bezig met opstarten")
		return
	}
	select {
	case <-time.After(a.werkduur):
		a.antwoord(w, http.StatusOK, "klaar")
	case <-r.Context().Done():
		// Verbinding verbroken, bijvoorbeeld doordat SIGKILL ons afkapt.
	}
}

// schakel zet een vlag om en meldt de nieuwe stand. Hetzelfde commando maakt
// dus stuk en repareert weer; handig als je live aan het typen bent.
func (a *app) schakel(vlag *atomic.Bool, naam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nieuw := !vlag.Load()
		vlag.Store(nieuw)
		log.Printf("debug: %s staat nu op %t", naam, nieuw)
		a.antwoord(w, http.StatusOK, fmt.Sprintf("%s=%t", naam, nieuw))
	}
}

// antwoord schrijft een statuscode met een leesbare regel erbij. Kubernetes
// kijkt alleen naar de code; de tekst is voor de zaal.
func (a *app) antwoord(w http.ResponseWriter, code int, tekst string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s: %s\n", a.pod, tekst)
}

func main() {
	log.SetFlags(log.Ltime)

	poort := tekstUitOmgeving("PORT", "8080")
	opstartduur := duurUitOmgeving("STARTUP_DELAY", 5*time.Second)

	naam, err := os.Hostname()
	if err != nil {
		naam = "onbekend"
	}
	a := &app{
		pod:      naam,
		werkduur: duurUitOmgeving("WORK_DURATION", 100*time.Millisecond),
	}

	// De server luistert meteen, maar alle endpoints antwoorden 503 tot het
	// opstarten klaar is. Zo krijgt de kubelet ook tijdens het opstarten een
	// nette statuscode in plaats van een geweigerde verbinding.
	go func() {
		log.Printf("opstarten: dit duurt %s", opstartduur)
		time.Sleep(opstartduur)
		a.opgestart.Store(true)
		log.Printf("opgestart: klaar voor verkeer")
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /startupz", a.startupz)
	mux.HandleFunc("GET /healthz", a.healthz)
	mux.HandleFunc("GET /readyz", a.readyz)
	mux.HandleFunc("GET /work", a.work)
	mux.HandleFunc("POST /debug/unready", a.schakel(&a.unready, "unready"))
	mux.HandleFunc("POST /debug/deadlock", a.schakel(&a.deadlock, "deadlock"))
	mux.HandleFunc("POST /debug/dependency-down", a.schakel(&a.dependencyDown, "dependency-down"))

	server := &http.Server{
		Addr:              ":" + poort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Bij SIGTERM zetten we readiness op 503 en maken we de lopende verzoeken
	// af. Let op: we wachten hier bewust niet zelf. Uit de EndpointSlice gaan
	// is de taak van de preStop-hook in het manifest - en juist die halen we in
	// scenario 5 weg. Zou de app zelf pauzeren, dan viel er niets te demonstreren.
	stopsignaal := make(chan os.Signal, 1)
	signal.Notify(stopsignaal, syscall.SIGTERM, os.Interrupt)

	afgesloten := make(chan struct{})
	go func() {
		defer close(afgesloten)
		<-stopsignaal
		a.afsluiten.Store(true)
		log.Printf("SIGTERM: readiness op 503, lopende verzoeken afmaken")

		ctx, annuleer := context.WithTimeout(context.Background(), afsluitTimeout)
		defer annuleer()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("afsluiten duurde te lang: %v", err)
		}
	}()

	log.Printf("luistert op :%s", poort)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server gestopt: %v", err)
	}
	<-afgesloten
	log.Printf("afgesloten")
}

// tekstUitOmgeving leest een omgevingsvariabele, met een standaardwaarde.
func tekstUitOmgeving(naam, standaard string) string {
	if waarde := os.Getenv(naam); waarde != "" {
		return waarde
	}
	return standaard
}

// duurUitOmgeving leest een duur als "90s" of "500ms". Onleesbare waarden
// laten de app niet omvallen; dat zou je tijdens een demo alleen maar ophouden.
func duurUitOmgeving(naam string, standaard time.Duration) time.Duration {
	waarde := os.Getenv(naam)
	if waarde == "" {
		return standaard
	}
	duur, err := time.ParseDuration(waarde)
	if err != nil {
		log.Printf("%s=%q is geen geldige duur, ik gebruik %s", naam, waarde, standaard)
		return standaard
	}
	return duur
}
