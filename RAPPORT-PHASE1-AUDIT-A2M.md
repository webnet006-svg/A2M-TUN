# RAPPORT PHASE 1 — AUDIT COMPLET (A2M TUN)

- **Date** : 15 novembre 2026
- **Branche** : `dams/ci` (parent `8cc3e845`, sous-module `clients/android` → `9117404`)
- **Méthode** : audit **statique** (lecture du code Kotlin/Compose, manifeste, navigation, ressources).
  Aucune modification de code dans cette phase.

---

## 1. Point d'entrée et architecture existante

| Élément | Fichier | Détail |
|---|---|---|
| Entrée app | `AndroidManifest.xml` | `MainActivity`, `MAIN`/`LAUNCHER`, `singleTask` |
| Contenu Compose | `MainActivity.kt:256-258` | `setContent { App() }` |
| Racine UI | `MainActivity.kt:392` | `fun App()` — Scaffold + NavHost + barres |
| Navigation | `Navigation.kt:112-116` | `NavHost`, `startDestination = "am_home"` (corrigé) |
| Bottom bar | `MainActivity.kt:1311-1326` | A2M bar si `currentRoute == "am_home"`, sinon barre Sing-box |
| Lecture état | `DashboardViewModel` (`uiState`) | état réel du service (`Status`) |

## 2. Ce qui affiche encore l'interface Sing-box

### 2.1 Écrans accessibles depuis l'interface A2M (routes Sing-box)
1. **Bottom bar A2M** (`A2MNavigationBar.kt:24-30`) :
   - Accueil → `am_home` (A2M)
   - **Serveurs → `dashboard`** (écran Sing-box : Profiles, groupes, connexions)
   - **Journal → `log`** (écran Logs Sing-box)
   - **Paramètres → `settings`** (écran Settings Sing-box, badge de mise à jour inclus)
2. **Menu hamburger A2M** (`A2MTUNHomeScreen.kt:278-305`) : même cartographie → mêmes écrans Sing-box.
3. **Carte serveur de l'accueil A2M** (`A2MTUNHomeScreen.kt:235-237`) :
   `dashboardViewModel.showProfilePickerSheet()` → **feuille de sélection de profils Sing-box**.
4. **Barre de statut service** (`MainActivity.kt:1071`) : `ServiceStatusBar` (composant Sing-box)
   visible en bas de l'écran A2M lorsque le service tourne ou démarre.

### 2.2 Navigation Sing-box complète toujours enregistrée
`Navigation.kt` déclare encore **~40 destinations** Sing-box, toutes atteignables par navigation :
- `dashboard`, `log`, `groups`, `connections`, `connections/detail/{id}`
- `profile/new`, `profile/edit/{profileId}`
- `tools` + `tools/network_quality`, `tools/stun_test`, `tools/usbip/*`, `tools/outbound_picker`,
  `tools/openconnect`, `tools/openvpn`, `tools/tailscale/*`, `tools/crash_report*`,
  `tools/oom_report*`, `tools/power_report*`
- `settings` + `settings/app`, `settings/core`, `settings/service`, `settings/profile_override`,
  `settings/profile_override/manage`, `settings/remote_control*`, `settings/privilege*`,
  `settings/tailscale/*`, `settings/fdroid_mirror`, `settings/privilege/logs`

### 2.3 Fenêtres de mise à jour GitHub (Sing-box)
| Élément | Localisation |
|---|---|
| Sollicitation 1er lancement (« check update ») | `MainActivity.kt:642-688` (`Settings.updateCheckPrompted`) |
| Dialog « mise à jour disponible » | `MainActivity.kt:692-724` (`UpdateAvailableDialog`, `UpdateState.updateInfo`) |
| Dialog de téléchargement/progression | `MainActivity.kt:726-780` (`UpdateState.downloadProgress`) |
| Badge mise à jour sur l'onglet Paramètres | `MainActivity.kt:1252-1262`, `1310` (`hasUpdate`) |
| Vérification réseau | `MainActivity.kt:242-249`, `670-673` (`Vendor.checkUpdateAsync()`) |

### 2.4 Éléments Sing-box visibles restants
- **Barre supérieure Sing-box** : titre « Dashboard » + menu 3 points (TopAppBar des écrans Sing-box)
  via `TopBarController` (`MainActivity.kt:507`, `1018-1020`, `1245`).
- **Bottom bar Sing-box** : `bottomNavigationScreens` = Dashboard, Logs, Tools, Settings
  (`NavigationDestinations.kt:56-63`), affichée pour toute route ≠ `am_home` (`MainActivity.kt:1327-1370`).
- **FAB de démarrage** Sing-box (`MainActivity.kt:1085-1200`, tablette) et barres `ServiceStatusBar` / `RemoteStatusBar`.
- **Nom de l'application** : `app_name = "sing-box"` (`res/values/strings.xml:4`).
- **ShareActivity** (`ShareActivity.kt:75`) : UI Sing-box lors d'un partage externe.

## 3. Moteur VPN à préserver (aucune suppression)

| Composant | Fichier |
|---|---|
| VpnService Android | `bg/VPNService.kt:20` |
| Gestion cœur + commandes | `bg/BoxService.kt:51` (`start/stop/serviceReload/closeService`) |
| Binder/Connexion service | `bg/ServiceBinder.kt`, `bg/ServiceConnection.kt` |
| Notification service | `bg/ServiceNotification.kt` |
| Tuile QS | `bg/TileService.kt` |
| État UI/commandes | `compose/screen/dashboard/DashboardViewModel.kt` (`toggleService:240`, `uiState:50`) |
| Logs réels | `compose/screen/log/LogViewModel.kt` (`clearLogs:104`, filtres, recherche) |
| Profils + import | `ProfilesCard.kt`, `ProfileImportHandler.kt`, `ClipboardImporter` |
| Statistiques | `DashboardUiState` (`downlink/uplink`, totaux, `serviceStartTime`) |

## 4. État actuel vs objectif

| Exigence mission | État |
|---|---|
| Démarrage direct sur A2M TUN | ✅ fait (`startDestination = "am_home"`, build vert `43665561404`) |
| Plus d'accès aux écrans Sing-box | ❌ bottom bar + menu + carte serveur y mènent encore |
| Plus de fenêtres de mise à jour GitHub | ❌ présentes |
| Identité nommée A2M TUN | ❌ `app_name` = « sing-box » |
| Pages Serveurs / Routage / Réglages A2M | ❌ inexistantes (routes Sing-box utilisées) |
| Journal A2M | ❌ écran Logs Sing-box |

## 5. Conclusion de l'audit

Le moteur est intact et exploitable. Le problème est **100 % couche applicative** :
la navigation A2M est branchée sur des **routes Sing-box** et l'ancienne UI reste
enregistrée dans le graphe de navigation, avec en plus la logique de mise à jour GitHub
et le nom « sing-box ».
