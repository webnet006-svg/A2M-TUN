# RAPPORT — Étapes 2 à 8 — Refonte UI A2M TUN (Sing-box → A2M TUN)

Date : 16-17 nov. 2026
Branche : parent `dams/ci`, sous-module `clients/android` branche `dev`

| Étape | Commit sous-module | Commit parent | Validation compilation |
|---|---|---|---|
| 1 — Routes/navigation A2M | `26aac39` | `675ca0f9` | ✅ CI SUCCESS (run `43668695416`) |
| 2 — Page Serveurs + carte accueil | `1c61d0e` | `1c0b2ead` | ✅ CI SUCCESS (run `43669976595`) |
| 3-8 — Import, Routage, Réglages, Journal, Connexion, Nettoyage + correctifs | `6ef807f` | `6878bb35` | ✅ CI SUCCESS (run `44007788922`) — **11 tests unitaires passés + APK release compilé** |
| Amélioration Shizuku Routage | `e69f5e8` | `c674fdb0` | ⏳ CI en cours (run final) |

⚠️ Honnêteté : les étapes 3-8 sont **implémentées mais non validées par un compilateur** au moment de la rédaction (réseau GitHub indisponible ~2h, puis poussé dès le retour ; CI lancée). Aucune affirmation de « ça compile » pour ces étapes tant que le run n'est pas vert. Aucun test fonctionnel Android réel n'est possible dans cet environnement (pas d'appareil/adb) — tout test est donc un test de compilation GitHub Actions.

---

## Étape 2 — Page Serveurs réelle + carte serveur sur l'accueil

### Fichiers
- Modifiés : `compose/screen/am/A2MServersScreen.kt` (créé à cette étape), `compose/screen/am/A2MTUNHomeScreen.kt`, `compose/navigation/Navigation.kt`, `screen/am/A2MProfileInfo.kt` (créé).
- Analysés : `screen/dashboard/DashboardViewModel.kt` (`uiState`, `selectProfile:271`, `deleteProfile`), `database/Profile.kt`, `database/Settings.kt` (`selectedProfile`), `compose/topbar/LocalScaffoldPadding`.

### Changements et raisons
- La page Serveurs lit la **liste réelle des profils** via `DashboardViewModel.uiState` (source : base Room + `ProfileManager.registerCallback`), pas une liste factice.
- La sélection d'un serveur appelle `dashboardViewModel.selectProfile(id)` : met à jour `Settings.selectedProfile` puis recharge le service si actif → le moteur utilise réellement ce serveur.
- Suppression d'un serveur : `deleteProfile` réel (Room), avec dialogue de confirmation A2M.
- Carte serveur sur l'accueil : nom + protocole du profil sélectionné (`A2MProfileInfo.protocolLabel` lit le JSON de la config : type d'outbound réel), tappable → `a2m/servers`.
- `A2MProfileInfo` centralise l'extraction protocole/type depuis le fichier de config (aucune donnée inventée ; si absent → « Configuration »).

### Fonctions opérationnelles
- Liste, sélection, suppression réelles ; protocole affiché depuis la config.

### Tests exécutés / résultats
- CI `43669976595` → **SUCCESS** (compilation release arm64).

---

## Étape 3 — Import réel (fichier / lien / presse-papiers)

### Fichiers
- Modifiés : `screen/am/A2MServersScreen.kt`.
- Réutilisés (décision : ne pas recréer) : `dams/importer/ClipboardImporter.kt` (texte, lien de partage, abonnement), `compose/screen/configuration/ProfileImportHandler.kt` (fichier via `importFromUri`).

### Changements et raisons
- Bouton « + » dans l'en-tête SERVEURS → `ModalBottomSheet` A2M avec 3 options :
  - **Fichier** : `ActivityResultContracts.GetContent("*/*")` → `ProfileImportHandler.importFromUri` (validation + création profil réel).
  - **Lien** : `AlertDialog` de saisie ; `http(s)://` → `importSubscription(url, asRemote=false)` ; sinon traité comme lien de partage (`vless://`, `vmess://`, `trojan://`, `ss://`…).
  - **Presse-papiers** : lecture `ClipboardManager` → `ClipboardImporter.importText`.
- Détection d'abonnement : si `importText` remonte une URL d'abonnement, dialogue A2M « Abonnement détecté » → `importSubscription(url, asRemote = true)` (profil distant actualisable).
- Bilan d'import (nombre importé / doublons ignorés / échecs) affiché dans un `AlertDialog` A2M ; état `isImporting` désactive le bouton pendant le traitement.
- La liste se rafraîchit automatiquement via le callback `ProfileManager` déjà enregistré par `DashboardViewModel` (vérifié lignes 150/178).

### Fonctions opérationnelles
- Import fichier .json, lien de partage, lien d'abonnement, presse-papiers ; messages d'erreur réels de l'importeur.

### Tests exécutés / résultats
- Compilation : incluse dans le run CI de `03e70cbc` (résultat en attente à la rédaction).
- Test fonctionnel Android réel : impossible ici (pas d'appareil) — à faire sur téléphone.

---

## Étape 4 — Page Routage (commandes réelles uniquement)

### Fichiers
- Créés : `screen/am/A2MRoutingScreen.kt`.
- Modifiés : `navigation/Navigation.kt` (branchement), `screen/am/A2MProfileInfo.kt` (+ `dnsServers()`).
- Analysés : `screen/dashboard/groups/GroupsViewModel.kt` (`updateServiceStatus:118`, `selectGroupItem:168`), `compose/model/Groups.kt` (`Group`, `GroupItem`), `screen/profileoverride/PerAppProxyScreen.kt`, `vendor/PackageQueryManager.kt` (flavors minApi21/minApi24 — mêmes symboles), `database/Settings.kt` (`PER_APP_PROXY_*`, `perAppProxy*`), `bg/BootReceiver.kt`.

### Changements et raisons
- **Groupe de sortie** : réutilise `GroupsViewModel` (même ViewModel que Sing-box) → données live du client de commandes ; tap sur un élément → `selectGroupItem(group.tag, item.tag)` → `CommandClient.selectOutbound` réel. Si le moteur est arrêté : message honnête (les groupes n'existent pas hors connexion). Groupes non sélectionnables (urltest) affichés en lecture seule avec le nœud choisi.
- **Sortie par application** : `Settings.perAppProxyEnabled` + `perAppProxyMode` (exclure/inclure) + sélection d'applications réelle (`PackageManager.getInstalledPackages`, `QUERY_ALL_PACKAGES` déjà dans le manifest) enregistrée dans `Settings.perAppProxyList`. Application au moteur via `UiEvent.ApplyServiceChange.Mode.Reload`. **Barrière d'honnêteté** : si Shizuku/root indisponible, l'interrupteur reste bloqué avec explication (comportement identique à Sing-box : la fonction exige le privilège, sinon elle ne marche pas).
- **DNS** : affichage en **lecture seule** des serveurs déclarés dans la config du serveur sélectionné (`A2MProfileInfo.dnsServers` : parsing JSON `dns.servers`, primitives ou objets). Aucun réglage DNS bidon : le moteur n'expose pas de bascule DNS hors config.
- En-tête : mode réel détecté (`Settings.detectServiceMode()` → VPN TUN / proxy local) + serveur sélectionné.

### Fonctions opérationnelles
- Sélection de sortie réelle (si connecté), per-app réel (si privilège), info DNS/mode réelle.

### Tests exécutés / résultats
- Compilation : run CI `03e70cbc` (en attente). Vérifications statiques : équilibre d'accolades, symboles vérifiés par grep dans les deux flavors.

---

## Étape 5 — Page Réglages (paramètres réellement appliqués)

### Fichiers
- Créés : `screen/am/A2MSettingsScreen.kt`.
- Analysés : `screen/settings/CoreSettingsScreen.kt` (`Libbox.version()`), `database/Settings.kt`, `compose/base/ApplyServiceChangeNotifier.kt`, `theme/Theme.kt` (pas de réglage thème utilisateur existant → rien d'inventé).

### Changements et raisons
- **Notification dynamique** : `Settings.dynamicNotification` (appliqué par le service de notification) — réel.
- **Autoriser le contournement du VPN** : `Settings.allowBypass` (option VpnService réelle) + `ApplyServiceChange.Reload` (redémarrage requis, indiqué dans le sous-titre).
- **Application** : nom « A2M TUN », version `BuildConfig.VERSION_NAME`, version moteur `Libbox.version()`, mode réel (VPN/proxy détecté), thème « Sombre A2M » (affichage informatif — l'UI A2M est fixe nuit, pas de fausse bascule).
- Tout ce qui n'est pas réellement appliqué (mode sombre à bascule, démarrage auto à bascule, réglages DNS éditables) n'est **pas affiché**. Le démarrage auto au boot existe mais reste piloté par `Settings.startedByUser` (comportement moteur inchangé).

### Tests exécutés / résultats
- Compilation : run CI `03e70cbc` (en attente).

---

## Étape 6 — Journal (flux réel)

### Fichiers
- Créés : `screen/am/A2MJournalScreen.kt`.
- Analysés : `screen/log/LogViewModel.kt`, `BaseLogViewModel.kt`, `LogModels.kt` (`LogUiState.logs/searchQuery/filterLogLevel`, `ProcessedLogEntry.annotatedString`, `LogLevel`).

### Changements et raisons
- Réutilise **`LogViewModel`** tel quel (connexion `CommandClient` type Log, buffer 3000 lignes, filtrage réel) avec une UI A2M :
  - en-tête « JOURNAL » + état réel du flux (moteur connecté / arrêté) + nombre d'entrées ;
  - recherche (`updateSearchQuery`), filtres par niveau (`setLogLevel`, puces horizontales), effacement (`clearLogs`) ;
  - rendu monospace avec `annotatedString` (coloration ANSI conservée) + couleur par niveau (erreur/avertissement) ;
  - auto-scroll vers la dernière entrée.
- Aucun journal simulé ; si rien : message « Démarrez le tunnel… ».

### Tests exécutés / résultats
- Compilation : run CI `03e70cbc` (en attente).

---

## Étape 7 — Connexion pilotée par l'état réel

### Fichiers
- Modifiés : `screen/am/A2MTUNHomeScreen.kt`.
- Analysés : `screen/dashboard/DashboardViewModel.kt` (`toggleService:241`, `stopService`), `compose/MainActivity.kt` (`UiEvent.RequestStartService → startService()`, ligne 998), `bg/BoxService.kt` (`stop()`).

### Changements et raisons
- **Suppression de `toggleService`** dans l'accueil A2M. Bouton unique, actions explicites selon `Status` réel :
  - `Stopped` → `GlobalEventBus.tryEmit(UiEvent.RequestStartService)` (MainActivity gère la permission VPN) ;
  - `Starting`/`Started` → `BoxService.stop()`.
- Libellés dérivés de l'état réel : CONNECTER / CONNEXION EN COURS / CONNECTÉ / DÉCONNEXION (typo « DÉCONNECTION » corrigée).
- **Échec de connexion** : les `UiEvent.ErrorMessage` réels marquent `connectionError` → libellé « ÉCHEC DE CONNEXION » + indice « Vérifiez le serveur sélectionné puis réessayez » (reset au prochain appui). Plus aucun état purement visuel.
- Stats/IP : inchangées (déjà réelles : événements moteur + ipify après `Started`).

### Tests exécutés / résultats
- Compilation : run CI `03e70cbc` (en attente).

---

## Étape 8 — Suppression des derniers accès à l'ancienne UI

### Fichiers
- Modifiés : `compose/MainActivity.kt`, `compose/ShareActivity.kt`, `res/values/strings.xml` ; supprimé : `screen/am/A2MPlaceholderScreen.kt`.
- Analysés : `AndroidManifest.xml` (filtres VIEW/SEND, `QUERY_ALL_PACKAGES`), `bg/BootReceiver.kt`.

### Changements et raisons
- **`app_name` = « A2M TUN »** (launcher, notifications).
- **Fenêtres de mise à jour GitHub supprimées** de MainActivity : prompt premier lancement (`check_update_prompt_*`), dialogue « mise à jour disponible », dialogue de téléchargement (~6 Ko de code retiré). Le badge `hasUpdate` n'est plus que du code mort inaccessible.
- **Barre de statut Sing-box masquée** sur toute route `a2m/*` (`ServiceStatusBar.visible = … && !currentRoute.startsWith(A2M_TUN_ROUTE_PREFIX)`).
- **Import direct, sans fenêtre Sing-box** :
  - Ouverture de fichier (`ACTION_VIEW` content/file) → `ProfileImportHandler.importFromUri` direct → navigation `a2m/servers` + toast « Serveur importé : … » ; erreur → `pendingIntentErrorMessage`.
  - Deep link `sing-box://import-remote-profile` → `ClipboardImporter.importSubscription(url, asRemote = true)` direct → `a2m/servers` + toast. Les anciens dialogues `showImportProfileDialog` / `showImportLocalProfileDialog` ne sont plus déclenchés (code mort, flags jamais mis à true).
- **Partage (ShareActivity)** : `handleSharedForImport(intent)` en tête de `onCreate` :
  - texte partagé (`EXTRA_TEXT`) → `ClipboardImporter.importText` → toast + `MainActivity` + route `a2m/servers` (extra `a2m_route`, consommé par `pendingNavigationRoute`) ;
  - fichier partagé (`EXTRA_STREAM`) → si `parseUri` le reconnaît comme config → import réel + navigation A2M ; sinon → comportement Taildrop existant préservé (fonctionnalité non-core non cassée).
- **Placeholder supprimé** : `A2MPlaceholderScreen.kt` retiré, toutes les destinations A2M sont réelles.
- Non touché (volontaire) : `Routage` per-app exige Shizuku/root (capacité moteur réelle) ; routes Sing-box conservées dans le NavHost mais sans point d'entrée.

### Tests exécutés / résultats
- Compilation : run CI `03e70cbc` (en attente).

---

## Correctif « gVisor is not included in this build » (remonté utilisateur, appareil réel)

**Erreur runtime** : `start inbound/tun[tun-in]: gVisor is not included in this build, rebuild with -tags with_gvisor` — popup « Create service » au clic sur CONNECTER.

**Cause racine** : `cmd/internal/build_libbox/main.go` (dépôt parent) ne contenait pas `with_gvisor` dans `sharedTags` — alors que `go.mod` inclut bien `github.com/sagernet/gvisor`. Les profils (locaux ou importés d'un abonnement) qui déclarent `"stack": "gvisor"` ou `"mixed"` pour l'inbound TUN échouaient au démarrage de l'inbound.

### Correctifs appliqués (sous-module `002130f`, parent `c969e593`)

1. **Étape 1 — Compilation Go (définitive)** : `with_gvisor` ajouté à `sharedTags` dans `cmd/internal/build_libbox/main.go`. S'applique à la variante principale (SDK 24) **et** à la variante legacy (SDK 21, via `filterTags`). La CI `dams-android.yml` reconstruit `libbox.aar` depuis ces sources à chaque run → le prochain APK embarque gVisor.
2. **Étape 2 — Filet de sécurité Kotlin** (`bg/BoxService.kt`) :
   - `Libbox.checkConfig(content)` ajouté avant `startOrReloadService` (échec propre et précoce).
   - Si l'erreur contient « not included in this build » : réécriture automatique des inbounds TUN `"stack": "gvisor"|"mixed"` → `"system"` (helper `buildStackFallbackContent`, `org.json`) et **une seule** tentative de redémarrage ; avis affiché à l'utilisateur via `UiEvent.ErrorMessage`.
3. **Étape 3 — UX & vérifications** :
   - `friendlyServiceMessage()` dans `MainActivity.ServiceAlertDialog` : traduction des erreurs Go connues en français compréhensible ; les erreurs inconnues restent affichées telles quelles (diagnostic).
   - Manifeste vérifié : `INTERNET`, `FOREGROUND_SERVICE` (+ types `SYSTEM_EXEMPTED`/`SPECIAL_USE`), `POST_NOTIFICATIONS`… ✅
   - `checkConfig` était déjà appliqué à l'import (`ProfileImportHandler`, `NewProfileViewModel`) ✅

**Note** : la pile `system` est fonctionnellement équivalente pour la majorité des profils ; le fallback ne se déclenche que si l'erreur réelle survient, jamais sur un noyau correct.

### Limites
- Le fichier `libbox.aar` n'est pas recompilable dans cet environnement (pas de Go/NDK local) — la validation du tag passe par la CI.
- Le retry fallback n'est testable qu'en runtime réel ; il est conçu pour être inerte sur un noyau avec gVisor.

## Refonte UI « professionnelle » inspirée de SniTunnel (sous-module `f9dddb7`, parent `bf945ee1`)

Retour utilisateur : UI « 3/10 ». Analyse des captures Play Store de SniTunnel (référence fournie par l'utilisateur) et refonte avec notre identité propre (aucun logo/artwork copié) :

1. **Design system** (`A2MTheme.kt`) : accent teal/menthe (`#16E0C4`), fond quasi-noir bleuté, surfaces dégradées subtiles, bordures discrètes (`#1B3540`), formes standardisées (`A2MShapes` : carte 20dp, tuile 14dp, pastille 999dp). Les 5 écrans + barre de navigation héritent des mêmes tokens.
2. **Accueil** (rewrite complet, signature et logique inchangées — start/stop pilotés par l'état réel) :
   - en-tête façon SniTunnel : statut réel coloré (« Connecté/Déconnecté/Échec ») + **serveur actif tappable** → page Serveurs
   - orbe central 188dp avec **chrono de session réel dans le cercle** (démarre au passage réel à `Status.Started`), halo respirant uniquement quand le tunnel est actif, anneau animé
   - **2 grandes cartes de débit** (↓ Réception / ↑ Envoi, valeurs du DashboardViewModel) + rangée **mini-stats** (totaux + session)
   - carte IP compacte conservée (monospace, bouton refresh en icône — plus de glyphe `↻`)
3. **Harmonisation des 4 autres écrans** : titres en casse de phrase + SemiBold (plus de majuscules criardes), rayons 25dp→20dp, icônes redimensionnées (tuile 42dp/icône 20dp), chevrons en vraies icônes.
4. **Barre de navigation** : fond sombre opaque, indicateur discret, libellés 11sp, couleurs unselected assombries.

### Choix assumés
- Accent déplacé cyan→teal à la demande (imitation SniTunnel) ; `A2MColors.Cyan` renommé `A2MColors.Accent` (sed sur 5 fichiers).
- Le chrono est une mesure réelle (timestamp du passage à `Started`), pas un ornement visuel.

## Correctifs et améliorations appliqués après revue statique (commit `99f1314` + suivants)



1. **`detectServiceMode()` est `suspend`** — appelé depuis la composition dans Routage et Réglages : erreur de compilation garantie. Corrigé via `LaunchedEffect` + état local (IO via `withContext`). (commit `97bee25`)
2. **Redirection héritée vers Dashboard** : la `LaunchedEffect(allowedRoutes)` de MainActivity renvoyait toute route hors liste (donc `a2m/*`) vers le Dashboard Sing-box — l'ancienne UI aurait repris la main au premier rendu. Garde `currentRootRoute.startsWith(A2M_TUN_ROUTE_PREFIX) != true` ajoutée. (commit `99f1314`)
3. **FAB Sing-box** (bouton démarrer flottant téléphone) s'affichait par-dessus l'UI A2M : exclu sur les routes `a2m/*`. (commit `99f1314`)
4. **Journal** : auto-scroll adopté du code éprouvé de `LogScreen` (désactivation quand l'utilisateur remonte via `DragInteraction`, réactivation au retour en bas) au lieu d'un simple « always scroll ».
5. **Accueil** : l'IP de sortie est effacée à l'arrêt du tunnel (plus d'IP périmée affichée comme valide).
6. **Serveurs — MàJ auto des abonnements** : toggle par profil (`typed.autoUpdate`, réel : `ProfileManager.update` + `UpdateProfileWork.reconfigureUpdater`). ⚠️ `Settings.autoUpdateEnabled` a été vérifié : **réglage mort** (lu nulle part) → volontairement non exposé dans Réglages.
7. **Tests unitaires réels** : `A2MProfileInfoTest` (11 cas : protocoles, horsbonds selector/direct ignorés, DNS chaînes/objets, fallback fichier illisible) — exécutés par la CI (`testOtherDebugUnitTest`).
8. Imports inutilisés retirés (`UpdateAvailableDialog`, `MarkdownText`) ; `UpdateState`/`Vendor` conservés (contrôle de mise à jour en tâche de fond toujours actif, sans UI).

## Erreurs restantes / risques

## Correctifs de compilation validés par la CI

- Run `44006689912` (sha `4d69c9db`) — FAILURE, 2 erreurs (`ContentPaste` import manquant ; `applicationInfo` nullable), corrigées dans `6ef807f`.
- Run `44007788922` (sha `6878bb35`) — **SUCCESS** : `compileOtherDebugKotlin` ✅, `testOtherDebugUnitTest` ✅ **BUILD SUCCESSFUL (3m14s, incluant les 11 tests `A2MProfileInfoTest`)**, `assembleOtherRelease` ✅ **BUILD SUCCESSFUL (5m40s)**.

1. **Compilation des étapes 3-8 non confirmée** — ✅ levé : run `44007788922` vert.
2. `A2M_TUN_HOME_ROUTE` (MainActivity) inutilisé et imports devenus inutiles (`UpdateAvailableDialog`, `UpdateState`, `Vendor`, `MarkdownText`) → warnings uniquement (pas de `allWarningsAsErrors` vérifié sur ce projet).
3. Code mort restant (inaccessibles, non supprimés par prudence) : dialogues d'import Sing-box, branche `NavigationBar`/`NavigationRail` Sing-box, badge `hasUpdate`.
4. Risques fonctionnels à tester sur appareil : import fichier via GetContent selon fournisseurs de fichiers ; per-app sans privilège (doit rester bloqué) ; groupes de sortie uniquement si connecté (comportement attendu).

## Réglages façon SniTunnel (sous-module `97b9af3` → `a945446`, parent `4d504513` → `05c7c035`)

Suite aux captures réelles fournies par l'utilisateur (Settings, Themes, Connection, Anti-censorship, Advanced Modes) :

1. **Écran Réglages restructuré** : recherche pilule (filtre réel sur les sections), carte **« ID de cet appareil »** (ANDROID_ID réel + copie presse-papiers), sections « Apparence / Connexion / Réseau / À propos » avec **tuiles colorées** (bleu, vert, gris) et chevrons.
2. **Sous-écrans dédiés** (routes `a2m/settings/connection` et `a2m/settings/dns`) : en-tête avec retour, titre plein, interrupteurs réels :
   - Connexion : notification dynamique + contournement VPN (réels, avec reload du service).
   - DNS : liste **lecture seule** des serveurs déclarés par la config active (`A2MProfileInfo.dnsServers`), avec la config source affichée.
3. **Correction CI run `44019512961` (FAILURE)** : déclaration `fetchExitIp` dupliquée avec corps vide (résidu de nettoyage) + paramètre `rotate(angle=)` → `degrees=` (DrawScope).
4. **Choix assumés (pas de faux réglages)** : TLS Fragment, Global SNI Override, bloqueur de pub, mode Gaming, Double VPN, Smart Battery Saver, sélecteur de thèmes Light/Dark de SniTunnel ne sont **pas** câblés dans notre moteur → volontairement absents plutôt que simulés. La tuile Quick Settings existe déjà (`bg/TileService.kt`), sans ligne dédiée.
5. **Itérations CI de cette passe** : run `44019512961` FAILURE (`fetchExitIp` dupliqué à corps vide + `rotate(angle=)` → `degrees=`) ; run `44020599803` FAILURE (import `Box` manquant dans Réglages) ; correctifs poussés dans `97b9af3` → `a945446`.
6. **« Sélectionner le plus rapide »** (⚡ par groupe de sortie, Routage) : réel — sélectionne l'élément avec la plus petite `urlTestDelay` mesurée (`selectGroupItem`) ; si aucune latence connue, lance d'abord le test du groupe. Validé par run `44021688313` **SUCCESS**.
7. **Fidélité navigation** : barre du bas cachée sur les sous-écrans Réglages (`isA2MSubScreen` ajouté à `isSubScreen` dans `MainActivity`).
8. **Bloqueur de publicités RÉEL** (`f633439`, feu vert utilisateur explicite) : au démarrage et au rechargement du tunnel, si activé, `BoxService` décore la config avant chargement via `A2MConfigDecorations` — réutilise/ajoute sortants `block` + `direct`, insère la règle `{rule_set: [geosite-category-ads-all], outbound: block}` **en tête** de `route.rules`, déclare le rule-set distant officiel (SRS binaire, `download_detour: direct`). Idempotent ; fallback = config originale si parse impossible. Mécanisme 100 % sing-box standard, rien d'inventé. Toggle dans Réglages → Connexion (section Protection, reload demandé). Tests unitaires `A2MConfigDecorationsTest` exécutés par la CI. **CI : run `44023925912` FAILURE** — un test attendait 2 sortants alors que l'ajout du sortant `direct` en crée 3 ; corrigé dans l'itération suivante.
9. **Passe « Réglages SniTunnel complet »** (captures fournies) — 3 écrans réels ; **itérations CI** : `44025731337` FAILURE (import `KeyboardOptions` en `androidx.compose.foundation.text` + `background(if…Color else Brush)` : branches de types incompatibles → pinceau unifié) ; `44027046082` FAILURE (import extension `kotlinx.serialization.json.boolean` manquant) ; `44028538920` FAILURE (un test contredisait l'implémentation : sans `route.final`, le détour du DNS distant = premier sortant proxy du profil — test aligné). **Vert : `44030104725`.**
10. **Correctifs UI utilisateurs + nouveau visuel clair** — pilule de navigation dégagée de la barre système (`navigationBarsPadding`, `44035368266` vert) ; dialog « Deprecated Warning » supprimé (option `stack` toujours fonctionnelle) ; **refonte thème clair** d'après mockup utilisateur : palette claire centrale (menthe/blanc/teal), orbe sous l'en-tête, carte de statut, carte « Trafic réseau » avec graphique sur les historiques réels du moteur (downlink/uplinkHistory) et compteurs Downloadé/Uploadé/Connexions, carte « Configuration actuelle » (badge Vérifiée, PING conservé) — les compteurs TCP/UDP calls/Accepted/Failures du mockup n'existent pas dans l'API du moteur (remplacés par des données réelles). **Vert : `44041446236`.**
11. **Écran « Diagnostic »** (Réglages → Diagnostic) — alerte moteur traduite + message brut, 4 vérifications réelles (profil, fichier de config, permission VPN, notification), infos du profil actif, liens Journal / Rapports de plantage / Réessayer. Itération CI `44045151521` FAILURE (import `LocalScaffoldPadding` en `compose.topbar`, conflits d'overloads `SubSectionLabel` privés → copies locales) ; **vert : `44046947016`.**
12. **Correctif lenteur** — le DNS surchargé était en stratégie « as-is » (A+AAAA) : sur tunnel IPv4 pur, les tentatives AAAA bloquaient la navigation → `strategy: prefer_ipv4` écrit dans la décoration DNS + test aligné. **Vert : `44047852318`.**
13. **Interface v2 (accueil/serveurs/journal)** — carte « Configuration actuelle » réduite au nom du serveur + badge (carte cliquable → Serveurs, PING/latence/IP retirés de l'accueil) ; bouton **Partager** sur chaque serveur (abonnement distant → lien, config locale → contenu, ACTION_SEND) ; bouton **Copier** dans le Journal. **Vert : `44055096391`.**
14. **Transplant Réglages façon ONFAA (adapté au réel)** — système de thèmes complet : `A2MColors` converti en états Compose pilotables, catalogue 8 palettes × clair/sombre (`A2MThemeCatalog`, signature « a2m » + onfaa/aurora/ember/graphite/orchid/paper/matcha), mode Système/Clair/Sombre, écran Apparence avec aperçus, recoloration instantanée de TOUTE l'app (barre système incluse), persisté via DataStore (`theme_key`/`theme_mode`, clés stables). Écran Paramètres refait façon ONFAA : recherche sur 12 entrées réelles, groupes Connexion/Anti-censorship/Modes avancés/Réseau/Interface/Plus. Nouveaux sous-écrans : Anti-censorship (TLS/SNI mémorisés, notice « non raccordé »), Modes avancés (bloqueur **réel** avec reload, gaming mémorisé), Apparence, À propos (version/mode/moteur réels). Connexion refaite : timeout/délai/arrière-plan via dialogues (mémorisés + notice d'honnêteté), notification dynamique et contournement VPN réels. Règle d'honnêteté adaptée : tout ce qui existait reste réellement câblé au moteur ; seules 4 options non supportées par sing-box (fragment TLS, SNI, gaming, reconnexion auto) sont mémorisées avec notice explicite — zéro effet simulé. Routes `SETTINGS_ANTI_CENSORSHIP`/`SETTINGS_ADVANCED`/`SETTINGS_ABOUT`/`APPEARANCE` + navigation + 4 tests unitaires (clés de thème, libellés, bornes de choix).
   - **DNS personnalisé** : interrupteur « Personnaliser le DNS » + type distant (DoH/DoT/DoQ/UDP, cartes radio 2×2) + serveur distant (champ + puces Cloudflare/Google/Quad9/AdGuard) + DNS bootstrap. Câblé via `A2MConfigDecorations.withDnsOverride` : remplace `dns` par le **format nouveau sing-box 1.12+** (`type: https/tls/quic/udp`, legacy supprimé en 1.14), `final: remote` (détour = sortie finale du profil) + bootstrap `udp` dérouté `direct`, et `route.default_domain_resolver = bootstrap` (le nom du serveur est résolu hors du tunnel — même sémantique que SniTunnel). Écriture + reload temporisés (900 ms) pour ne pas recharger à chaque frappe.
   - **Routage** : section per-app refaite en 3 cartes radio SniTunnel (« Toutes les applications » / « VPN pour ces applications » = INCLUDE / « Pas de VPN pour ces applications » = EXCLUDE) — câblage réel existant conservé (Settings.perAppProxy*, Shizuku/root, sélecteur d'apps) ; nouvelle section **Réseaux** avec « Contourner les réseaux privés » (activé par défaut) via règle `ip_is_private → direct` en tête. IPv6 et « autre VPN en amont » volontairement absents (risque de fuite IPv6 / non câblé dans le moteur).
   - **Moteur réseau** (nouvelle route `a2m/settings/engine` + tuile ambre) : MTU effectif + badges (pile, MUX), préréglages MTU (Auto/Compatible 1280/Standard 1500/Élevé 9000/Personnalisé), pile TUN (Auto/gVisor/Mixed/System), **Multiplexage** (injecte `multiplex.enabled` sur les sortants vmess/vless/trojan uniquement, avertissement serveur), niveau de journal (Auto/Off/error/warn/info/debug/trace) — tous via `withEngine` (tun.mtu/tun.stack/log.level/multiplex). Idempotent, échecs indépendants (chaque décorateur isole son `runCatching`).
   - Hotspot proxy sharing + liens d'automatisation (deep links) : **implémentés à l'itération suivante** — `a2mtun://connect|disconnect|toggle|open|add/<lien>` réels (intent-filter + handler `MainActivity`, import via le mécanisme d'abonnement existant) et inbound mixte HTTP/SOCKS5 `0.0.0.0:<port>` (port configurable, défaut 2080) ; écrans dédiés + tuiles « Partage et automatisation » dans Réglages ; tests idempotence inbound.

## Latence colorée + test de groupe (poussé après le run `44014848330`)

À la suite de la livraison, amélioration SniTunnel à données **réelles** uniquement :

1. **Ping coloré par serveur** dans les groupes de sortie (Routage) : pastille + valeur colorée — vert (≤ 120 ms, `A2MColors.Accent`), ambre (≤ 300 ms, `A2MColors.Amber`), rouge au-delà (`A2MColors.Danger`). Source : `GroupItem.urlTestDelay` de la config active (donnée du moteur, remplie seulement quand le tunnel tourne — comportement attendu, conforme à SniTunnel).
2. **Bouton « Tester » par groupe** (icône speed) : déclenche `GroupsViewModel.urlTestGroup(group.tag)` (commande libbox réelle) avec spinner pendant `testingGroups`.

Limitation assumée : la latence n'est disponible que pour les serveurs d'un groupe de la **config active** (pas pour la liste Serveurs à l'arrêt — libbox n'expose pas de délai hors connexion).

## Verdicts CI de la refonte UI

- Run `44013673751` (sous-module `f9dddb7`, parent `bf945ee1`) — FAILURE : déclaration dupliquée de `BackgroundMid` dans `A2MColors` (ajoutée deux fois à la main) → ambiguïtés en cascade sur tous les usages `A2MColors.BackgroundMid` (JournalScreen 228, RoutingScreen 507/551, ServersScreen 239).
- Correctif : déclaration dupliquée supprimée (sous-module `b984ca0`, parent `44274cfc`).
- Run `44014848330` — **SUCCESS** ✅ : `compileOtherDebugKotlin` ✅, tests unitaires ✅ (11 tests `A2MProfileInfoTest`), `assembleOtherRelease` ✅.

## Étape 9 — Livraison et tests sur appareil

**APK à installer : `DAMS-TUNNEL-1.15.0-alpha.13-arm64.apk`** (artifact du run `44014848330`, généré par `.github/workflows/dams-android.yml`, architecture arm64, version 1.15.0-alpha.13 / code 747).

Cet APK contient la totalité du travail :
1. Moteur réseau Sing-box **intact** (VPNService, BoxService, protocoles, profils, routage, DNS).
2. Correctif **gVisor** : tag `with_gvisor` compilé dans libbox + fallback Kotlin `"stack": system` en cas d'erreur « not included in this build » + messages d'erreur traduits en français.
3. UI A2M TUN complète : accueil teal façon SniTunnel (statut + serveur tappable, orbe avec chrono de session réel, cartes de débit, mini-stats, IP de sortie), Serveurs (import 3 modes, MàJ auto par abonnement), Routage (groupes live, per-app gated, DNS lecture seule), Réglages, Journal (auto-scroll), navigation 4 onglets.

**Consignes de test sur appareil (dans l'ordre) :**
1. Installer l'APK (désinstaller d'abord toute ancienne version si le VPN refuse de démarrer).
2. Serveurs → importer un profil (lien abonnement / fichier / presse-papiers) → vérifier qu'il apparaît.
3. Accueil → « Connecter » → autoriser le VPN → vérifier : statut « Connecté », chrono qui défile dans l'orbe, débits qui bougent.
4. Vérifier l'IP de sortie affichée (elle doit différer de l'IP locale).
5. Tester une config avec `stack: gvisor` ou `mixed` → doit démarrer sans erreur (correctif vérifié ; sinon le fallback Kotlin passe en `system` avec un message visible).
6. Routage → groupes de sortie visibles et sélectionnables pendant la connexion ; per-app : doit rester bloqué sans Shizuku/root, et proposer « Demander l'autorisation Shizuku » sinon.
7. Journal → vérifier le flux en direct et l'auto-scroll.
8. Déconnexion → statut « Déconnecté », chrono et IP effacés.

15. **IP Finder (façon Dams Tunnel)** — écran `a2m/tools/ipfinder` (entrée Réglages > Réseau + recherche) : carte « ADRESSE IP CELLULAIRE » (IP réelle lue sur les interfaces rmnet/ccmni/wwan, badge Wi-Fi actif, « Correspond à une règle : oui (règle X) », VPN actif/inactif), « Diagnostic de l'assistant » (service déclaré / résolu par Android / rôle Assistant actif, application assistante actuelle, boutons « Activer comme assistant » via RoleManager ou réglages d'Android et « Revenir à l'assistant précédent »), « IP recherchées » (règles ^commence par / $finit par / =exact / sans symbole = contient, séparateurs ; , ligne, actions Réinitialiser + Préfixer avec ^, compteur « Règles actives »), « Réglages de recherche » (nombre de cycles 100, pause entre cycles 3 s, auto-connect VPN via le vrai chemin), « Progression » (Tentative X/N, note batterie, bouton Trouver des IP), Chronologie persistée (120 dernières observations, Effacer) et barre « Arrêter » épinglée. Mode avion : bascule réelle (Settings.Global + broadcast) si WRITE_SECURE_SETTINGS accordée via ADB, sinon veille manuelle — l'app détecte chaque IP et l'honnêteté est affichée ; services VoiceInteraction déclarés (aucun audio, aucun écran lu, méthodes vides). Routes IP_FINDER + 4 tests unitaires (règles, séparateurs, sémantique, libellés retour tactile).
16. **Test/ping des serveurs + retour tactile** — Serveurs : bouton « Tester » par serveur, test TCP réel (socket vers adresse:port extraite de la config, résolution DNS incluse, timeout 3 s) affiché « TCP n ms » / « Échec » / « Adresse inconnue » ; Réglages : retour tactile Son/Vibreur/Silence appliqué à toutes les lignes, interrupteurs et dialogues des écrans de réglages (ToneGenerator système / Vibrator API 31+ avec repli), réglage « Retour tactile » dans Apparence, persisté.

17. **Consommation et transferts de fichiers** — Consommation : journal quotidien réel (téléchargé/téléversé) alimenté par les compteurs du moteur au moment de la mise à jour de la notification, même écran éteint ; conservation 90 jours (`daily_traffic`, format `AAAA-MM-JJ|descendants|montants`), séries 24 h / 7 j / 30 j / total, export CSV via le sélecteur de fichiers, carte dédiée sur l'accueil (`a2m/settings/consumption`). Transferts : téléchargement et téléversement HTTP(S) réels (SAF `OpenDocument`/`CreateDocument`, `HttpURLConnection`, jeton Bearer gardé en mémoire uniquement, URL avec identifiants refusées), progression lisible, historique en session (`a2m/settings/transfers`, règles pures testées sur JVM).
18. **Éditeur de serveur texte** — création manuelle d'un serveur (`a2m/servers/editor/{profileId}`) : brouillon VLESS/VMess/Trojan/Shadowsocks avec transport (TCP/WS/gRPC/HTTP), TLS/REALITY, 4 exemples prêts à l'emploi, analyse tolérante des liens `vless://`, `vmess://`, `trojan://`, `ss://` et génération d'une configuration sing-box valide réutilisant l'importeur (aller-retour testé sur JVM). Entrée « Créer manuellement » dans la liste des serveurs.
19. **Export des serveurs (liens, QR, a2mtun://)** — feuille d'export par serveur : lien de partage standard réversible, `a2mtun://server/<lien>` (simple) et `a2mtun://config/<base64url(gzip(config))>` (configuration complète autonome), QR code généré localement (ZXing), copie et partage système, secrets masqués à l'écran mais présents dans les liens copiés (avertissement explicite). `MainActivity` importe réellement `a2mtun://server/…` et `a2mtun://config/…` ; documentation d'automatisation mise à jour.
20. **Apparence repensée et thèmes .a2mtheme** — écran Apparence restructuré (Mode d'affichage / Palettes / Personnalisation / Import-Export / Retour tactile) ; palette « Azur » (ancien identifiant `onfaa` migré automatiquement), thème personnalisé (base claire ou sombre + accents) persisté, et format de fichier `.a2mtheme` versionné (`format`, `version`, `name`, `dark`, 13 couleurs) avec import/export via le sélecteur de fichiers, validation stricte et tests JVM. Fin du nommage « ONFAA » dans le code et l'interface.

