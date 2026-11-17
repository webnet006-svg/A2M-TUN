# Rapport de test — Intégration interface A2M TUN (client Android)

- **Date** : 15 novembre 2026
- **Branche** : `dams/ci`
- **HEAD** : `09ae9dd5` (parent DAMS-TUNNEL) — sous-module `clients/android` → `97207dc` (branche `dev`)
- **Version testée** : `1.15.0-alpha.13`
- **Outil de test** : GitHub Actions (`.github/workflows/dams-android.yml`), job « Build arm64 release APK »

---

## Contexte

Phase 2 de l’intégration de l’interface HTML/CSS/JS A2M TUN
(commits `7963b27` / `15ea9239`) dans le client Android
(`A2MTUNHomeScreen.kt`, `A2MNavigationBar.kt`, wiring navigation).

La CI de release compile le client Android via la tâche Gradle
`:app:compileOtherDebugKotlin` ; c’est sur ce chemin que les erreurs
ont été observées.

## Résultats des runs CI

| Run | Commit | Conclusion | Erreurs constatées | Correction |
|---|---|---|---|---|
| `43657498516` | `15ea9239` | **failure** | `A2MNavigationBar.kt:45` Icon ; `A2MTUNHomeScreen` collectAsState / mutableStateOf ; clip / border / fillMaxHeight ; selectedProfile… ; fetchExitIp suspend ; let inference | `2d5c974`, `97207dc` |
| `43659196517` | `1abc46ba` | **failure** | `A2MTUNHomeScreen:7` clip (package `androidx.compose.foundation` incorrect) ; `fetchExitIp` appel dans lambda `onRefresh` non-suspend | `97207dc` |
| `43659973075` | `09ae9dd5` | **✅ success** | — | — |

## Corrections apportées (run final, `43659973075`)

1. `A2MNavigationBar.kt` : ajout de l’import `androidx.compose.material3.Icon`.
2. `A2MTUNHomeScreen.kt` : ajout des imports manquants :
   - `androidx.compose.runtime.collectAsState`
   - `androidx.compose.runtime.mutableStateOf`
   - `androidx.compose.foundation.border`
   - `androidx.compose.foundation.layout.fillMaxHeight`
   - `androidx.compose.ui.draw.clip` *(package correct pour `Modifier.clip`)*.
3. `fetchExitIp` déplacée en **top-level suspend fun** ;
   rafraîchissement IP via `rememberCoroutineScope().launch { … }`.

## Résultat final

- ✅ **Run `43659973075` : SUCCESS** (durée 12 min 15 s).
- ✅ `Run unit tests` : OK (suite existante, sans nouvelle régression).
- ✅ `Build release APK` : OK.
- ✅ Artefact publié : `dams-android-arm64` (34,2 Mo).
- ✅ APK release téléchargé : `/tmp/opencode/apk/DAMS-TUNNEL-1.15.0-alpha.13-arm64.apk`.

## Correctif final — écran de démarrage (run `43665561404`, SUCCESS)

- **Problème constaté** : l’application affichait toujours l’écran Sing-box
  (`Screen.Dashboard`) au démarrage, l’écran A2M TUN (`am_home`) n’étant jamais
  affiché initialement.
- **Cause** : `startDestination = Screen.Dashboard.route` dans le `NavHost`
  (`Navigation.kt`).
- **Correctif** : `startDestination = "am_home"`.
- **Validation** : run `43665561404` → **SUCCESS** (build release APK OK,
  unit tests OK).

## Artefacts

- Release APK (intégration A2M TUN + start destination fix, build final) :
  `/tmp/opencode/apk/DAMS-TUNNEL-1.15.0-alpha.13-arm64.apk`
- Debug APK (build antérieur à l’intégration A2M, run `43546409135`) :
  `/tmp/opencode/apk-debug/other/debug/SFA-1.15.0-alpha.13-*-debug.apk` +
  `otherLegacy`

## Blocs / non-testés

- Tests terrain sur téléphone réel (checklist étape 0) : à réaliser.
- Tests unitaires applicatifs : suite existante (21 tests) — aucune nouvelle régression observée.
