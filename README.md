# Leeft je app nog?

Workshopmateriaal over probes en health checks in Kubernetes, voor beginners.
De repo bevat drie dingen:

- **`presentatie/slides.md`** — de presentatie (21 slides, [Marp](https://marp.app/))
- **`demo/`** — een Go-app die je op commando kunt stukmaken, plus het image
- **`demo/k8s/`** — het goede manifest en vier varianten met precies één fout

De slides staan online op **https://miguelmartens.github.io/leeft-je-app-nog/**.
Het draaiboek voor de live demo staat in [`DRAAIBOEK.md`](DRAAIBOEK.md).

## Wat je nodig hebt

| Tool                                               | Waarvoor                            |
| -------------------------------------------------- | ----------------------------------- |
| [kind](https://kind.sigs.k8s.io/)                  | het lokale cluster                  |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | met het cluster praten              |
| [Docker](https://docs.docker.com/get-docker/)      | het image bouwen                    |
| [Go](https://go.dev/dl/) 1.27                      | alleen als je de app wilt aanpassen |

Kubernetes 1.30 of nieuwer, want het manifest gebruikt de `sleep`-preStop-hook.
De node-image die kind standaard meelevert voldoet.

Gebruik je Podman in plaats van Docker:

```bash
make build DOCKER=podman
export KIND_EXPERIMENTAL_PROVIDER=podman
```

## Zelf draaien

```bash
make alles      # kind-cluster + image bouwen + laden + uitrollen
```

Dat is hetzelfde als:

```bash
make cluster    # 1 control plane, 2 workers
make build      # image als probe-demo:dev
make load       # image in het kind-cluster
make deploy     # het goede manifest
```

Controleren of het werkt:

```bash
curl http://127.0.0.1:30080/work
kubectl get pods -o wide
```

`make help` laat alle commando's zien. Opruimen doe je met `make schoon`.

## De demo-app

Eén Go-bestand, geen dependencies, met de endpoints uit de presentatie:

```
GET  /startupz               ben ik opgestart?
GET  /healthz                leef ik nog?          (startup + liveness probe)
GET  /readyz                 mag ik verkeer?       (readiness probe)
GET  /work                   het echte werk

POST /debug/unready          readiness aan/uit
POST /debug/deadlock         /healthz laten hangen, liveness gaat falen
POST /debug/dependency-down  de nagebootste database aan/uit
```

De `/debug`-endpoints zijn schakelaars: hetzelfde verzoek maakt stuk en
repareert weer. Ze zitten niet achter de Service, dus je bereikt ze via een
port-forward naar één specifiek pod:

```bash
kubectl port-forward pod/<naam> 8080:8080
curl -X POST http://127.0.0.1:8080/debug/unready
```

`make unready` en `make dependency` doen dit voor je.

Instelbaar via de omgeving: `PORT` (8080), `STARTUP_DELAY` (5s) en
`WORK_DURATION` (100ms).

## De vier kapotte varianten

Elk bestand in `demo/k8s/kapot/` is een volledig manifest dat op precies één
punt afwijkt van `demo/k8s/deployment.yaml`. Bovenin staat wat er kapot is en
wat de fix is.

| Bestand                             | Kapot                            | Gevolg                              |
| ----------------------------------- | -------------------------------- | ----------------------------------- |
| `1-zonder-readiness.yaml`           | geen readinessProbe              | verkeer naar pods die nog opstarten |
| `2-liveness-checkt-dependency.yaml` | liveness wijst naar `/readyz`    | alle replica's herstarten tegelijk  |
| `3-trage-start.yaml`                | 90s opstarten, geen startupProbe | CrashLoopBackOff                    |
| `4-zonder-prestop.yaml`             | geen preStop, grace period 1s    | gedropte requests bij een update    |

Het verschil bekijken:

```bash
diff -u demo/k8s/deployment.yaml demo/k8s/kapot/1-zonder-readiness.yaml
```

Uitrollen met `make kapot-1` t/m `make kapot-4` — die tonen eerst de diff.
Terug naar goed met `make herstel`.

## De slides

Snelste weg: installeer de extensie **Marp for VS Code** en open
`presentatie/slides.md`. Je krijgt live preview en presentatiemodus.

Zonder VS Code:

```bash
make slides     # live preview in je browser

npx @marp-team/marp-cli@latest presentatie/slides.md --html -o dist/index.html
npx @marp-team/marp-cli@latest presentatie/slides.md --pdf --pdf-notes -o dist/slides.pdf
```

Druk in de HTML op `p` voor de presenter view met sprekersnotities, de volgende
slide en een timer. De notities staan in het markdown-bestand als
HTML-commentaar onder elke slide.

Alle opmaak zit in een `<style>`-blok bovenin `slides.md` — geen losse
theme-bestanden, geen build-config. De kleuren staan er als CSS-variabelen;
de drie probe-kleuren (amber voor startup, rood voor liveness, groen voor
readiness) lopen door de hele presentatie heen.

### Publiceren

`.github/workflows/pages.yaml` bouwt de slides bij elke push naar `main` en zet
ze op GitHub Pages. Zet daarvoor eenmalig **Settings → Pages → Source** op
_GitHub Actions_.
