---
marp: true
theme: default
paginate: true
html: true
title: Leeft je app nog?
description: Een introductie tot probes in Kubernetes
author: Miguel Martens
---

<style>
:root {
  --ink: #1B2027;
  --ink2: #2B333D;
  --soft: #F1F3F5;
  --muted: #5C6672;
  --muted-d: #9AA6B2;
  --light: #D7DDE4;
  --green: #1FA97A;
  --amber: #D98C1F;
  --red: #C4364C;
  /* lichtere varianten, alleen voor tekst op een donkere achtergrond */
  --green-l: #6FD9AE;
  --amber-l: #F0B45C;
}

section {
  font-family: "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  font-size: 24px;
  color: var(--ink);
  background: #FFFFFF;
  padding: 60px 70px;
}

h1 {
  font-family: Georgia, "Times New Roman", serif;
  font-size: 52px;
  margin: 0 0 6px 0;
  color: var(--ink);
}

h2 {
  font-family: Georgia, "Times New Roman", serif;
  font-size: 30px;
  margin: 0 0 10px 0;
}

/* de ondertitel onder elke H1 */
h1 + p em {
  display: block;
  font-style: normal;
  font-weight: 600;
  font-size: 16px;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: var(--muted);
  margin-bottom: 28px;
}

code { font-size: 0.85em; }

pre {
  background: var(--ink2);
  border-radius: 10px;
  padding: 20px 24px;
  font-size: 20px;
  line-height: 1.45;
}
pre code { color: var(--light); background: none; }

/* Syntaxkleuring. Marp levert de GitHub-kleuren voor een lichte achtergrond:
   een string wordt #0a3069 en een yaml-key #0550ae. Dat donkerblauw valt weg
   tegen het codeblok, en op een beamer lees je er helemaal niets meer van.
   Hieronder dezelfde rolverdeling in de kleuren van deze deck: wit voor de
   structuur, groen voor de waarde, amber voor het sleutelwoord. */
pre {
  --color-prettylights-syntax-storage-modifier-import: var(--light);
  --color-prettylights-syntax-comment: var(--muted-d);
  --color-prettylights-syntax-constant: #FFFFFF;
  --color-prettylights-syntax-entity: #FFFFFF;
  --color-prettylights-syntax-markup-heading: #FFFFFF;
  --color-prettylights-syntax-markup-bold: #FFFFFF;
  --color-prettylights-syntax-markup-italic: var(--light);
  --color-prettylights-syntax-string: var(--green-l);
  --color-prettylights-syntax-string-regexp: var(--green-l);
  --color-prettylights-syntax-entity-tag: var(--green-l);
  --color-prettylights-syntax-constant-other-reference-link: var(--green-l);
  --color-prettylights-syntax-keyword: var(--amber-l);
  --color-prettylights-syntax-variable: var(--amber-l);
  --color-prettylights-syntax-markup-list: var(--amber-l);
}

/* kolommen */
.cols  { display: grid; grid-template-columns: 1fr 1fr; gap: 26px; }
.cols3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 22px; }
.cols-wide { display: grid; grid-template-columns: 3fr 2fr; gap: 26px; }

.card {
  background: var(--soft);
  border-radius: 10px;
  padding: 24px 26px;
  font-size: 21px;
}
.card h2 { font-size: 26px; }
.card.dark { background: var(--ink2); color: var(--light); }
.card.dark h2 { color: #FFFFFF; }
.card ul { margin: 0; padding-left: 20px; }
.card li { margin-bottom: 10px; }
.card p { color: var(--muted); }
.card.dark p, .card.dark li { color: var(--light); }

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px; height: 38px;
  border-radius: 50%;
  color: #fff; font-weight: 700; font-size: 19px;
  margin-right: 12px;
  vertical-align: middle;
}
.s { background: var(--amber); }
.l { background: var(--red); }
.r { background: var(--green); }
.n { background: var(--ink); }

.note {
  font-style: italic;
  color: var(--muted);
  font-size: 21px;
  margin-top: 26px;
}

.rows { display: grid; gap: 12px; }
.row {
  background: var(--soft);
  border-radius: 10px;
  padding: 16px 22px;
  font-size: 21px;
  display: grid;
  grid-template-columns: 46px 1fr;
  align-items: center;
}
.row.wide { grid-template-columns: 46px 1fr 1fr; gap: 16px; }
.row .dim { color: var(--muted); font-size: 19px; }

/* donkere slides */
section.dark {
  background: var(--ink);
  color: #FFFFFF;
}
section.dark h1, section.dark h2 { color: #FFFFFF; }
section.dark h1 + p em { color: var(--amber); }
section.dark .note { color: var(--muted-d); }
section.dark .row { background: var(--ink2); color: var(--light); }

/* titelslide */
section.title { background: var(--ink); color: #FFFFFF; justify-content: center; }
section.title h1 { color: #FFFFFF; font-size: 68px; }
section.title .sub { font-style: italic; font-size: 30px; color: var(--light); margin-top: 4px; }
section.title .legend { margin-top: 46px; font-size: 22px; color: var(--light); }
section.title .legend span.badge { margin-left: 34px; }
section.title .legend span.badge:first-child { margin-left: 0; }
section.title .who { margin-top: 70px; font-size: 19px; color: var(--muted-d); }

/* slotslide met de repo-link */
section.title .repo { margin-top: 46px; font-size: 26px; }
section.title .repo a {
  color: #FFFFFF;
  text-decoration: none;
  border-bottom: 2px solid var(--amber);
  padding-bottom: 4px;
}

/* Het complete manifest: 23 regels yaml naast 22, en dat past niet op 20px.
   Alleen op deze slide krimpt de code en levert de kop wat hoogte in, zodat
   het hele manifest in beeld blijft in plaats van onderaan af te lopen. */
section.manifest h1 { font-size: 42px; }
section.manifest h1 + p em { margin-bottom: 14px; }
section.manifest pre {
  font-size: 15px;
  line-height: 1.32;
  padding: 14px 18px;
}

section.title::after, section.dark::after { color: var(--muted-d); }
</style>

<!-- _class: title -->
<!-- _paginate: false -->

# Leeft je app nog?

<div class="sub">Een introductie tot probes in Kubernetes</div>

<div class="legend">
<span class="badge s">S</span> startup
<span class="badge l">L</span> liveness
<span class="badge r">R</span> readiness
</div>

<div class="who">Miguel Martens</div>

<!--
Openingsvraag: wie heeft ooit een pod gehad die Running was terwijl de applicatie
niet werkte? Bijna alle handen gaan omhoog. Daar gaat deze sessie over.
-->

---

# Wat je na dit uur weet

*Geen voorkennis nodig, alleen een cluster en nieuwsgierigheid*

<div class="rows">
<div class="row"><span class="badge r">✓</span> Waarom Kubernetes niet zelf kan zien of je applicatie werkt</div>
<div class="row"><span class="badge r">✓</span> Wat de drie probes doen en wanneer je welke gebruikt</div>
<div class="row"><span class="badge r">✓</span> Waar je ze neerzet in je YAML, en wat je app moet antwoorden</div>
<div class="row"><span class="badge r">✓</span> De vier fouten die vrijwel iedereen de eerste keer maakt</div>
</div>

<div class="note">En daarna gaan we het live stukmaken. Dat is het leukste deel.</div>

<!--
Zet de verwachting: dit is een introductie. Aan het eind heeft iedereen zelf
een probe kapotgemaakt en gerepareerd.
-->

---

# Draaien is niet hetzelfde als werken

*Het probleem dat probes oplossen*

<div class="cols">
<div class="card">

## Wat Kubernetes standaard ziet

Het proces in de container draait. Meer niet. Dat is de hele controle.

```
NAME          READY   STATUS
web-7d9f...   1/1     Running
```

</div>
<div class="card dark">

## Wat er ondertussen aan de hand kan zijn

- De app is nog bezig met opstarten
- De app hangt en beantwoordt niets meer
- De database is even onbereikbaar
- De app draait, maar geeft op elk verzoek een fout

</div>
</div>

<div class="note">Kubernetes kan niet in je applicatie kijken. Jij moet het vertellen — en dat doe je met probes.</div>

<!--
Dit is de kernboodschap van de hele sessie. Neem er de tijd voor;
een praktijkverhaal werkt hier goed.
-->

---

# Wie stelt die vraag eigenlijk?

*De kubelet — het onderdeel van Kubernetes dat op elke node draait*

<div class="cols3">
<div class="card">
<span class="badge n">1</span>

## De kubelet vraagt

Op elke node draait een kubelet. Die beheert de containers op die node en stuurt periodiek een verzoek naar jouw container.

</div>
<div class="card">
<span class="badge n">2</span>

## Jouw app antwoordt

Meestal met een HTTP-statuscode. 200 betekent goed, 503 betekent niet goed. Jij bepaalt in je code wat goed betekent.

</div>
<div class="card">
<span class="badge n">3</span>

## Kubernetes handelt

Op basis van het antwoord: de container herstarten, of geen verkeer meer naar dit pod sturen. Automatisch.

</div>
</div>

<div class="note">Belangrijk: de kubelet gaat rechtstreeks naar je container. Niet via een Service, niet via je ingress, niet via DNS.</div>

<!--
Als iemand vraagt wat een kubelet is: de agent die op elke worker node draait
en daar de containers beheert. Meer hoeven ze nu niet te weten.
-->

---

# Drie probes, drie vragen

*Elke probe stelt precies één vraag*

<div class="cols3">
<div class="card">
<span class="badge s">S</span> **Startup**

*"Ben je al opgestart?"*

Draait alleen tijdens het opstarten. Zolang deze loopt, laten de andere twee je met rust.

</div>
<div class="card">
<span class="badge l">L</span> **Liveness**

*"Leef je nog?"*

Merkt op dat je app vastzit. Het proces draait wel, maar er komt niets meer uit. Antwoord: herstarten.

</div>
<div class="card">
<span class="badge r">R</span> **Readiness**

*"Mag je verkeer aan?"*

Draait de hele tijd door. Bepaalt of gebruikers bij dit pod terechtkomen. Geen herstart.

</div>
</div>

<div class="note">Stel je geen probe in? Dan gaat Kubernetes er altijd van uit dat het antwoord goed is.</div>

<!--
Liveness en readiness zorgen voor de meeste verwarring.
Dat lossen we op de volgende slide op.
-->

---

# En als het antwoord 'nee' is?

*Twee heel verschillende gevolgen — dit is het belangrijkste onderscheid*

<div class="cols">
<div class="card">

<h2><span class="badge l">L</span> Liveness faalt</h2>

**Kubernetes herstart je container.**

Het proces wordt gestopt en opnieuw gestart. Alles in het geheugen is weg. Je pod blijft bestaan, maar telt een herstart op in de RESTARTS-kolom.

</div>
<div class="card">

<h2><span class="badge r">R</span> Readiness faalt</h2>

**Kubernetes stopt met verkeer sturen.**

Je pod wordt uit de Service gehaald — de load balancer voor je pods. Gebruikers komen bij de andere replica's terecht. Je container blijft draaien.

</div>
</div>

<div class="note">Vuistregel: liveness kijkt naar binnen (doe ik het nog?), readiness kijkt naar buiten (kan ik iemand helpen?).</div>

<!--
Als er één slide blijft hangen, dan deze.
Herhaal de vuistregel later nog eens bij de fouten-slide.
-->

---

# Waar zet je dit neer?

*In de containerdefinitie van je pod, naast image en ports*

<div class="cols-wide">
<div>

```yaml
spec:
  containers:
    - name: app
      image: mijn-app:1.0
      ports:
        - containerPort: 8080

      livenessProbe:          # <- hier
        httpGet:
          path: /healthz
          port: 8080

      readinessProbe:         # <- en hier
        httpGet:
          path: /readyz
          port: 8080
```

</div>
<div class="card">

## Let op

- Per container, niet per pod — twee containers hebben elk hun eigen probes
- Verschillende paden voor liveness en readiness. Dit is geen detail
- Het pad is gewoon een URL in je eigen applicatie, die je zelf schrijft

</div>
</div>

<!--
Laat dit ook live zien in hun eigen Helm values of met kubectl edit.
Beginners willen weten waar het fysiek staat.
-->

---

# Hoe stelt de kubelet de vraag?

*Vier manieren — begin met de eerste*

<div class="cols">
<div class="card dark">

## `httpGet` — begin hier

Een HTTP-verzoek naar een pad in je app. Goed bij een antwoord tussen 200 en 399. Dit is wat je meestal wilt.

</div>
<div class="card">

## `tcpSocket`

Kan er verbinding gemaakt worden met de poort? Zegt niets over de applicatie erachter. Voor dingen die geen HTTP spreken.

</div>
<div class="card">

## `grpc`

Voor gRPC-services, met het standaard health-protocol van gRPC. Alleen relevant als je app gRPC gebruikt.

</div>
<div class="card">

## `exec`

Een commando in de container. Klinkt handig, maar start elke keer een nieuw proces — dat kost je nodes merkbaar CPU.

</div>
</div>

<div class="note">Twijfel je? Neem httpGet. Je hebt er alleen een endpoint in je app voor nodig.</div>

<!--
Beginners kiezen vaak exec met curl omdat dat bekend voelt.
Waarschuw ervoor — veel images hebben helemaal geen curl.
-->

---

# De knoppen en hun standaardwaarden

*Je hoeft ze niet allemaal in te vullen, maar ken ze*

| Veld | Standaard | Wat het doet |
|---|---|---|
| `periodSeconds` | 10 | Hoe vaak de kubelet het vraagt. |
| `timeoutSeconds` | 1 | Hoelang hij op antwoord wacht. Eén seconde is krap — zet hem gerust op 3. |
| `failureThreshold` | 3 | Hoe vaak het achter elkaar mis mag gaan voordat Kubernetes ingrijpt. |
| `initialDelaySeconds` | 0 | Wachttijd voor de allereerste vraag. Met een startup probe overbodig. |
| `successThreshold` | 1 | Hoe vaak het goed moet gaan om weer als gezond te tellen. |

<div class="note">Rekenvoorbeeld: elke 10 seconden vragen en 3 keer mogen falen betekent dat Kubernetes na ongeveer 30 seconden ingrijpt.</div>

<!--
timeoutSeconds op 1 seconde is de meest onderschatte default.
Bij een drukke app faalt dan de probe, niet de app.
-->

---

# Start je app langzaam op?

*Dan heb je een startup probe nodig*

<div class="cols">
<div class="card">

## Het probleem

Je app heeft 60 seconden nodig om op te starten. De liveness probe begint al na 10 seconden te vragen en grijpt na 30 seconden in.

Je app wordt dus herstart voordat hij ooit klaar is. Steeds opnieuw. Je ziet CrashLoopBackOff, en het lijkt een bug in je code.

</div>
<div>

**De oplossing** — een startup probe zet liveness en readiness op pauze tot het opstarten klaar is.

```yaml
startupProbe:
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 5
  failureThreshold: 24   # 120s de tijd
```

</div>
</div>

<div class="note">Zo krijgt je app rustig twee minuten om op te starten, terwijl de liveness probe daarna gewoon scherp staat.</div>

<!--
Dit is demo-scenario 4. Laat eerst de CrashLoopBackOff zien, laat ze raden
wat er mis is, en voeg dan pas de startup probe toe.
-->

---

# Wat moet je app antwoorden?

*Zo'n endpoint schrijven is minder werk dan je denkt*

<div class="cols-wide">
<div>

```go
// liveness: kijkt alleen naar zichzelf
http.HandleFunc("/healthz", func(w, r) {
    w.WriteHeader(200)
})

// readiness: mag wel dependencies checken
http.HandleFunc("/readyz", func(w, r) {
    if !db.Bereikbaar() {
        w.WriteHeader(503)
        return
    }
    w.WriteHeader(200)
})
```

</div>
<div class="card">

## De regels

- 200 is goed, 503 is niet goed
- Houd het snel: geen zware queries, geen schrijfacties
- Geen authenticatie op deze endpoints — de kubelet logt niet in
- Liveness checkt nooit een database. Readiness mag dat wel
- Maar: checken alle replica's dezelfde database, dan gaan ze ook allemaal tegelijk uit de Service

</div>
</div>

<!--
Benadruk hoe simpel liveness mag zijn: komt het verzoek binnen en wordt het
beantwoord, dan leeft het proces. Dat is genoeg.
-->

---

# Waarom heet het /healthz?

*Een stukje Google-geschiedenis dat in je YAML is blijven hangen*

<div class="cols">
<div class="card">

## Z-pages bij Google

Interne diensten bij Google kregen automatisch een setje diagnostische endpoints: `/varz` voor metrics, `/statusz` voor status, `/rpcz` voor verkeer. Samen heetten ze z-pages.

Die z zat erachter om botsingen te voorkomen met echte URL's in de applicatie zelf — een app had vaak al een eigen `/status`.

</div>
<div class="card">

## Waar je het nu nog ziet

- Kubernetes zelf: de API-server heeft `/livez` en `/readyz`
- Prometheus gebruikt hetzelfde patroon — geschreven door ex-Googlers
- In vrijwel elk voorbeeld op internet, inclusief deze presentatie
- Je mag ook gewoon `/health` nemen. Kubernetes kijkt alleen naar de statuscode

</div>
</div>

```console
$ kubectl get --raw='/livez?verbose'      # probeer dit eens op je eigen cluster
```

<!--
Leuk weetje om de aandacht even vast te houden na een blok techniek.
Draai het kubectl-commando live: de API-server laat dan zijn eigen checklist zien.
Benadruk dat het een gewoonte is en geen standaard - er is geen RFC die /healthz voorschrijft.
-->

---

# Een complete Dockerfile

*Je code is klaar — nu verpakken we hem in een image*

<div class="cols-wide">
<div>

```dockerfile
# ---------- 1. bouwen ----------
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app .

# ---------- 2. draaien ----------
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
```

</div>
<div class="card">

## Wat hier gebeurt

- Twee keer `FROM`: je bouwt in een image mét Go, je draait in een image zonder
- `go.mod` apart kopiëren — die laag verandert zelden, dus je builds blijven snel
- `CGO_ENABLED=0` maakt een statisch bestand, want distroless heeft geen libc
- Resultaat: klein image, geen shell, geen package manager, draait als nonroot

</div>
</div>

<div class="note">Valt je iets op? Er staat geen HEALTHCHECK in. Dat is met opzet.</div>

<!--
Loop de Dockerfile regel voor regel door. Voor veel deelnemers is multi-stage nieuw.
De EXPOSE-regel is puur documentatie: die opent niets, en Kubernetes doet er niets mee.
-->

---

# Twee verrassingen met je image

*Hier loopt iedereen een keer tegenaan*

<div class="cols">
<div class="card">

<h2><span class="badge l">1</span> HEALTHCHECK doet niets</h2>

Docker en Docker Compose kijken naar de HEALTHCHECK-instructie. Kubernetes negeert hem volledig en leest alleen de probes in je Pod-spec.

</div>
<div class="card">

<h2><span class="badge l">2</span> Er zit geen curl in je image</h2>

In een distroless of scratch image zit geen curl, geen wget en geen shell. Een exec-probe faalt daar meteen, met een foutmelding die nergens op slaat.

</div>
</div>

```yaml
# faalt: curl bestaat niet in dit image
livenessProbe:
  exec: { command: ["curl", "-f", "http://localhost:8080/healthz"] }

# werkt altijd: de kubelet doet het HTTP-verzoek zelf
livenessProbe:
  httpGet: { path: /healthz, port: 8080 }
```

<!--
Als iemand vraagt waarom hun exec-probe een exec format error geeft:
dit is het antwoord. Dit kost mensen soms een halve dag.
-->

---

<!-- _class: manifest -->

# Het complete manifest

*Van je code, via je image, naar een draaiende pod*

<div class="cols">
<div>

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mijn-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mijn-app
  template:
    metadata:
      labels:
        app: mijn-app
    spec:
      terminationGracePeriodSeconds: 30
      containers:
        - name: app
          image: mijn-app:1.0
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: 50m
```

</div>
<div>

```yaml
          # 1. ben ik opgestart?
          startupProbe:
            httpGet:
              path: /healthz
              port: 8080
            periodSeconds: 5
            failureThreshold: 24

          # 2. leef ik nog?
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            periodSeconds: 10
            timeoutSeconds: 3

          # 3. mag ik verkeer?
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            periodSeconds: 3
```

</div>
</div>

<!--
Laat zien dat de rechterkolom letterlijk doorloopt waar de linker ophoudt,
allemaal binnen dezelfde container. Wijs op de inspringing: probes staan op
hetzelfde niveau als image en ports. Vertel er ook bij dat je hier normaal nog
een lifecycle preStop en een Service onder zet - die staan in de repo.
-->

---

<!-- _class: dark -->

# Vier fouten die vrijwel iedereen maakt

*En die er in een code review prima uitzien*

<div class="rows">
<div class="row"><span class="badge l">✕</span><div><strong>Hetzelfde pad voor liveness en readiness</strong><br><span class="dim">Dan heb je geen readiness meer, maar een tweede manier om te herstarten.</span></div></div>
<div class="row"><span class="badge l">✕</span><div><strong>De liveness probe checkt de database</strong><br><span class="dim">De database hapert even, alle replica's falen tegelijk en herstarten tegelijk. Een klein probleem wordt een storing.</span></div></div>
<div class="row"><span class="badge l">✕</span><div><strong>De readiness probe checkt een gedeelde dependency</strong><br><span class="dim">Alle replica's gaan tegelijk uit de Service. Nul endpoints, terwijl de pods zelf prima werken - ook voor verzoeken die de database niet nodig hebben.</span></div></div>
<div class="row"><span class="badge l">✕</span><div><strong>Geen resources ingesteld op de container</strong><br><span class="dim">Een container met te weinig CPU kan zijn eigen probe niet op tijd beantwoorden en herstart zichzelf eindeloos.</span></div></div>
</div>

<div class="note">De tweede zie je straks live in de demo — dat is het meest indrukwekkende scenario.</div>

<!--
Fout 2 is de belangrijkste. Vraag de zaal wat er gebeurt als tien replica's
tegelijk herstarten op een database die het al zwaar heeft.
Fout 3 is de spiegel daarvan: niemand herstart, maar er blijft ook geen pod over
om verkeer naartoe te sturen. Uit de Service halen helpt alleen als er nog een
gezond pod is - bij een gedeelde database is dat er niet.
-->

---

# Netjes afsluiten

*De laatste fout: waarom je bij een update toch foutmeldingen ziet*

Als een pod verdwijnt bij een update, gebeuren er twee dingen tegelijk: je container krijgt het signaal om te stoppen, en Kubernetes haalt het pod uit de Service. Er is geen garantie wat het eerst klaar is. Zonder pauze sluit je app af terwijl er nog verzoeken onderweg zijn — en die worden foutmeldingen.

<div class="cols3">
<div class="card"><span class="badge n">1</span>

**Wacht even**
Een korte pauze voordat je afsluit, zodat Kubernetes je pod uit de Service kan halen

</div>
<div class="card"><span class="badge n">2</span>

**Vang het stopsignaal op**
Geen nieuwe verzoeken aannemen, de lopende netjes afmaken

</div>
<div class="card"><span class="badge n">3</span>

**Geef jezelf de tijd**
`terminationGracePeriodSeconds` ruim boven je pauze plus je langste verzoek

</div>
</div>

```yaml
lifecycle:
  preStop: { sleep: { seconds: 5 } }      # terminationGracePeriodSeconds: 30
```

<!--
Dit is demo-scenario 5 en meestal de grootste eyeopener: dezelfde app,
één regel YAML verschil, van tientallen fouten naar nul.
-->

---

# De demo-app

*Één klein Go-programma met de endpoints van zojuist*

<div class="cols">
<div>

```console
GET  /startupz    ben ik opgestart?
GET  /healthz     leef ik nog?
GET  /readyz      mag ik verkeer?
GET  /work        het echte werk

POST /debug/unready
POST /debug/deadlock
POST /debug/dependency-down
```

</div>
<div class="card">

## Hoe het werkt

- Met de debug-endpoints breek je de app met één commando
- Geen herbouwen of opnieuw uitrollen tijdens de sessie
- Draait op kind, dus je hebt geen cloudcluster nodig
- Je krijgt de repo mee om er zelf mee te spelen

</div>
</div>

<!--
Laat de code even zien: bewust één bestand, leesbaar op een projector.
-->

---

# Zo zie je wat er gebeurt

*Drie commando's die je daarna blijft gebruiken*

<div class="rows">
<div class="row"><span></span><div><code>kubectl get pods -w</code><br><span class="dim">De RESTARTS-kolom loopt op als je liveness probe faalt. READY gaat van 1/1 naar 0/1 bij readiness.</span></div></div>
<div class="row"><span></span><div><code>kubectl describe pod &lt;naam&gt;</code><br><span class="dim">Onderaan bij Events staat letterlijk waarom een probe faalde, inclusief de statuscode.</span></div></div>
<div class="row"><span></span><div><code>kubectl get endpointslice</code><br><span class="dim">Laat zien welke pods verkeer krijgen. Hier verdwijnt een pod zodra readiness faalt.</span></div></div>
</div>

<div class="note">Een EndpointSlice is simpelweg de lijst met pods waar een Service verkeer naartoe stuurt.</div>

<!--
Laat ze deze drie commando's zelf intikken voordat de demo begint.
Dan volgen ze straks alles mee.
-->

---

# Demo: vijf keer stuk, vijf keer gemaakt

*Eerst raden wat er gaat gebeuren, dan pas enter*

<div class="rows">
<div class="row wide"><span class="badge l">1</span><strong>Update zonder readiness probe</strong><span class="dim">gebruikers krijgen foutmeldingen tijdens een uitrol</span></div>
<div class="row wide"><span class="badge l">2</span><strong>Readiness uitzetten</strong><span class="dim">het pod verdwijnt uit de lijst, maar blijft draaien</span></div>
<div class="row wide"><span class="badge l">3</span><strong>Liveness die de database checkt</strong><span class="dim">alle pods herstarten tegelijk</span></div>
<div class="row wide"><span class="badge l">4</span><strong>Trage start zonder startup probe</strong><span class="dim">CrashLoopBackOff die op een bug lijkt</span></div>
<div class="row wide"><span class="badge r">5</span><strong>Netjes afsluiten</strong><span class="dim">van tientallen foutmeldingen naar nul</span></div>
</div>

<!--
Reken op vijf tot zeven minuten per scenario.
Vraag elke keer eerst aan de zaal wat ze verwachten.
-->

---

<!-- _class: dark -->

# Vier dingen om te onthouden

<div class="rows">
<div class="row"><span class="badge r">1</span> Liveness kijkt naar binnen, readiness kijkt naar buiten. Nooit omdraaien.</div>
<div class="row"><span class="badge l">2</span> Liveness checkt nooit een database of een andere service.</div>
<div class="row"><span class="badge s">3</span> Start je app langzaam op? Dan gebruik je een startup probe.</div>
<div class="row"><span class="badge r">4</span> Een probe die nog nooit gefaald heeft, is een probe die je niet getest hebt.</div>
</div>

<div class="note">Vragen? En daarna: naar de terminal.</div>

<!--
Vraag welk scenario ze morgen als eerste in hun eigen cluster gaan proberen.
De repo-link staat op de volgende slide.
-->

---

<!-- _class: title -->
<!-- _paginate: false -->

# Alles staat online

<div class="sub">Slides, demo-app en de vijf scenario's</div>

<div class="repo"><a href="https://github.com/miguelmartens/leeft-je-app-nog">github.com/miguelmartens/leeft-je-app-nog</a></div>

<div class="who">Miguel Martens</div>

<!--
Laat deze slide staan tijdens de vragen, zodat iedereen de link kan overtypen.
In de README staat hoe ze het cluster in tien minuten zelf opzetten.
-->
