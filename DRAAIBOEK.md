# Draaiboek

Voor mij tijdens de sessie. Scanbaar, niet om voor te lezen.
Reken op 5 tot 7 minuten per scenario.

---

## Opzetten (voor de zaal binnenkomt)

```bash
make alles          # cluster + image + het goede manifest
```

Vier terminals, groot lettertype:

| #   | Commando         | Wat je laat zien                                    |
| --- | ---------------- | --------------------------------------------------- |
| 1   | `make watch`     | `kubectl get pods -w` — READY en RESTARTS           |
| 2   | `make endpoints` | `kubectl get endpointslice -w` — wie krijgt verkeer |
| 3   | `make belasting` | verkeer naar `/work`, elke fout een regel           |
| 4   | vrij             | hier typ je de commando's hieronder                 |

Terminal 3 blijft de hele sessie doorlopen. Stil = goed, regels = fout.

> Kijk voor elk scenario eerst naar terminal 3. Nul fouten is het nulpunt.

---

## 1 — Update zonder readiness probe

**Doel:** verkeer gaat naar pods die nog niet klaar zijn.

```bash
make kapot-1                              # toont eerst de diff
kubectl rollout restart deployment/probe-demo
```

**Vraag vooraf:** _"Deze pods hebben geen readiness probe. Wat gebeurt er met de gebruikers als ik nu een nieuwe versie uitrol?"_

**Wat de zaal ziet:**

- Terminal 1: nieuwe pods gaan meteen naar `1/1 Running`
- Terminal 2: dat pod staat direct in de EndpointSlice
- Terminal 3: **`FOUT 503`** — vijf seconden per pod, en de app start nog op

**De uitleg:** zonder readiness probe zegt Kubernetes altijd "ja, klaar". Draaien is niet hetzelfde als werken.

**Fix:**

```bash
make herstel
kubectl rollout restart deployment/probe-demo   # nu wél zonder fouten
```

---

## 2 — Readiness live uitzetten

**Doel:** readiness haalt je uit de Service, maar herstart niets.

```bash
make unready        # POST /debug/unready naar één pod
```

**Vraag vooraf:** _"Ik zet zo bij één pod readiness op 503. Gaat die pod herstarten, ja of nee?"_

**Wat de zaal ziet:**

- Terminal 1: `READY` gaat van `1/1` naar `0/1`, **RESTARTS blijft 0**
- Terminal 2: het IP van dat pod verdwijnt uit de EndpointSlice
- Terminal 3: blijft stil — de andere twee replica's vangen alles op

**De uitleg:** liveness kijkt naar binnen, readiness kijkt naar buiten. Dit is readiness, dus geen herstart.

**Fix:**

```bash
make unready        # hetzelfde commando zet het weer aan
```

---

## 3 — Liveness die de dependency checkt

**Doel:** het anti-pattern. Eén trage dependency herstart je hele deployment.

```bash
make kapot-2                    # liveness wijst nu naar /readyz
kubectl rollout status deployment/probe-demo
make dependency                 # 'database' onbereikbaar op alle pods
```

**Vraag vooraf:** _"De liveness probe checkt hier de database, zoals in heel veel echte manifests. De database hapert twee minuten. Wat gebeurt er?"_

**Wat de zaal ziet:**

- Terminal 1: eerst `0/1`, dan na ~30 seconden **RESTARTS 1, 2, 3 — bij alle drie tegelijk**
- Terminal 3: alles fout, want er is geen enkele gezonde replica meer over
- `kubectl describe pod <naam>` → Events: `Liveness probe failed: HTTP probe failed with statuscode: 503`

**De uitleg:** een klein probleem is een storing geworden. Bij tien replica's op een database die het al zwaar heeft, is dit een sneeuwbal.

**Fix:**

```bash
make herstel        # liveness terug naar /healthz
make dependency     # dependency weer omhoog
```

Na `herstel` blijven de pods draaien met een falende readiness — laat dat even staan, dat is precies goed gedrag. Daarna `make dependency` en alles komt terug.

---

## 4 — Trage start zonder startup probe

**Doel:** CrashLoopBackOff die op een bug in je code lijkt.

```bash
make kapot-3        # STARTUP_DELAY 90s, geen startupProbe
```

**Vraag vooraf:** _"Deze app heeft 90 seconden nodig om op te starten. De liveness probe vraagt elke 10 seconden en mag 3 keer falen. Reken even mee."_

**Wat de zaal ziet:**

- Terminal 1: nieuwe pod, na ~30 seconden `RESTARTS 1`, daarna `2`, dan **`CrashLoopBackOff`**
- Het oude pod blijft staan (`maxUnavailable: 0`), dus terminal 3 blijft stil
- `kubectl describe pod <naam>` → `Liveness probe failed`, **niet** een crash van de app
- `kubectl logs <naam>` → `opstarten: dit duurt 1m30s` — de app doet niets fout

**De uitleg:** dit ziet eruit als een bug en is een YAML-probleem. Laat ze eerst raden voordat je de startup probe noemt.

**Fix:**

```bash
make herstel        # startupProbe terug: 24 x 5s = 120s de tijd
```

Wil je het overtuigender: rol `3-trage-start.yaml` mét startupProbe uit en laat zien dat dezelfde 90 seconden nu gewoon goed gaan.

---

## 5 — Netjes afsluiten

**Doel:** de eyeopener. Dezelfde app, één blok YAML verschil, van tientallen fouten naar nul.

```bash
make kapot-4                                    # geen preStop, grace period 1s
kubectl rollout restart deployment/probe-demo
```

**Vraag vooraf:** _"Alle probes staan hier goed. Toch gaat deze update fout. Waarom?"_

**Wat de zaal ziet:**

- Terminal 3: **een stroom fouten** (`000` en `502`/`503`) tijdens de uitrol
- Terminal 2: het pod-IP staat er nog in terwijl de container al weg is

**De uitleg:** afsluiten en uit de Service gehaald worden gebeuren tegelijk, zonder volgorde. `preStop` koopt kube-proxy de tijd om bij te werken.

**Fix — dit is de afsluiter:**

```bash
make herstel
kubectl rollout restart deployment/probe-demo   # zelfde update, nul fouten
```

Laat terminal 3 nog tien seconden stil staan voordat je iets zegt.

---

## Als er iets vastloopt

```bash
kubectl describe pod <naam>          # onderaan bij Events staat de statuscode
kubectl get pods -o wide             # over welke nodes staan ze verdeeld
make logs                            # logs van alle replica's tegelijk
make herstel                         # altijd terug naar het goede manifest
kubectl rollout restart deployment/probe-demo
```

Helemaal opnieuw: `make schoon && make alles` (ongeveer twee minuten).
