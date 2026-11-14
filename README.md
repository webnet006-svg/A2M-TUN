# DAMS-TUNNEL

**DAMS-TUNNEL** est un fork du projet [sing-box](https://github.com/SagerNet/sing-box) (« The universal proxy platform »), orienté vers la construction et la distribution d’un client Android (arm64-v8a) via GitHub Actions.

> Ce dépôt est un travail dérivé/community fork. Il n’est pas affilié au projet sing-box/SagerNet. Pour la documentation fonctionnelle complète de sing-box, référez-vous à la documentation officielle en amont.

## État du projet

- **Upstream :** [SagerNet/sing-box](https://github.com/SagerNet/sing-box)
- **Branche par défaut du dépôt :** `testing`
- **Branche de build Android/CI :** `dams/ci`
- **Cible actuelle :** APK Android **arm64-v8a** (release signée en CI à partir des secrets configurés dans le dépôt)

## Compilation Android (GitHub Actions)

Le workflow `.github/workflows/dams-android.yml` (branche `dams/ci`) :

- se déclenche sur les **push** vers `dams/ci` ;
- construit `libbox` pour Android/arm64 (Go, Android NDK, JDK 17) ;
- initialise le sous-module `clients/android` ;
- exécute les tests unitaires et génère un **APK release arm64** ;
- publie un artefact et une **pre-release** GitHub taguée `dams-v<VERSION>-arm64`.

**Artefacts / builds :**
- Exécutions Actions : [Actions › DAMS Android Release](https://github.com/webnet006-svg/DAMS-TUNNEL/actions/workflows/dams-android.yml)
- Pre-releases / APK : [Releases](https://github.com/webnet006-svg/DAMS-TUNNEL/releases)

> Remarque : la signature release et certaines propriétés de build dépendent de secrets GitHub (`KEYSTORE_BASE64`, `LOCAL_PROPERTIES`). Sans ces secrets, le job de build release ne peut pas aboutir.

## Build local (aperçu)

Ce projet hérite de la structure de sing-box (Go). Pour le client Android, la CI :
1. construit `libbox` (`make lib_install`, puis `go run ./cmd/internal/build_libbox -target android -platform android/arm64`) ;
2. prépare le client Android (`clients/android`) et produit l’APK via Gradle.

Les prérequis exacts (SDK/NDK, versions Go/JDK, configuration de signature) sont définis dans le workflow CI.

## Documentation upstream

- Documentation sing-box : https://sing-box.sagernet.org
- Dépôt upstream : https://github.com/SagerNet/sing-box

## Licence

DAMS-TUNNEL est basé sur sing-box et reste soumis aux conditions de licence du projet amont. Voir le fichier [`LICENSE`](LICENSE) pour le texte complet.

Avis de licence hérité (sing-box) :

```
Copyright (C) 2022 by nekohasekai <contact-sagernet@sekai.icu>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

In addition, no derivative work may use the name or imply association
with this application without prior consent.
```