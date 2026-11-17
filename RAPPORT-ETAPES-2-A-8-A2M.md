# RAPPORT — Étapes 2 à 8 — Refonte UI A2M TUN (Sing-box → A2M TUN)

Date : 16-17 nov. 2026
Branche : parent `dams/ci`, sous-module `clients/android` branche `dev`

| Étape | Commit sous-module | Commit parent | Validation compilation |
|---|---|---|---|
| 1 — Routes/navigation A2M | `26aac39` | `675ca0f9` | ✅ CI SUCCESS (run `43668695416`) |
| 2 — Page Serveurs + carte accueil | `1c61d0e` | `1c0b2ead` | ✅ CI SUCCESS (run `43669976595`) |
| 3-8 — Import, Routage, Réglages, Journal, Connexion, Nettoyage | `b9979ee` | `03e70cbc` | ❌ CI FAILURE (run `44006689912`) — 2 erreurs de compilation, corrigées |
| Correctifs compilation | `6ef807f` | `6878bb35` | ⏳ CI en cours |

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

### CI run `44006689912` (sha `4d69c9db`) — FAILURE, corrigé dans `6ef807f`
- `A2MRoutingScreen.kt:129` : `PackageInfo.applicationInfo` est **nullable** (API) — passage direct à `getApplicationLabel` refusé. Corrigé avec `?.let { … } ?: packageName`.
- `A2MServersScreen.kt:254` : `Icons.Filled.ContentPaste` non résolu — import explicite `androidx.compose.material.icons.filled.ContentPaste` manquant (les icônes étendues exigent un import par symbole). Ajouté.
- Ce sont les **2 seules erreurs** rapportées par `compileOtherDebugKotlin` ; le reste du module compilait.

1. **Compilation des étapes 3-8 non confirmée** — run CI de `03e70cbc` en cours ; toute erreur relevée sera corrigée avant l'étape 9.
2. `A2M_TUN_HOME_ROUTE` (MainActivity) inutilisé et imports devenus inutiles (`UpdateAvailableDialog`, `UpdateState`, `Vendor`, `MarkdownText`) → warnings uniquement (pas de `allWarningsAsErrors` vérifié sur ce projet).
3. Code mort restant (inaccessibles, non supprimés par prudence) : dialogues d'import Sing-box, branche `NavigationBar`/`NavigationRail` Sing-box, badge `hasUpdate`.
4. Risques fonctionnels à tester sur appareil : import fichier via GetContent selon fournisseurs de fichiers ; per-app sans privilège (doit rester bloqué) ; groupes de sortie uniquement si connecté (comportement attendu).

## Prochaine étape
- **Étape 9** : vérifier run CI `03e70cbc` → corriger si erreur → tests complets + APK release arm64 + rapport final.
