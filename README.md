# Groupie Tracker

## Auteurs
* TSCHAEN Morgane
* BOYCHEVA Ivana

Projet réalisé dans le cadre du cursus Ynov Campus Strasbourg - Année 2025-2026

## Objectif du projet

Groupie Tracker est une application web développée en Go permettant de visualiser, filtrer et explorer les données d'artistes musicaux et leurs concerts via l'API Groupie Trackers. L'application met l'accent sur une architecture backend robuste avec toute la logique côté serveur.

## Comment lancer le serveur

### Prérequis
* Go 1.21 ou supérieur

### Installation et lancement

```bash
# Cloner le repository
git clone https://github.com/m-tschaen/Groupie-Tracker_BOYCHEVA_TSCHAEN.git
cd Groupie-Tracker_BOYCHEVA_TSCHAEN

# Lancer le serveur
go run .
```

Le serveur démarre sur `http://localhost:8080`

## Routes principales

| Route | Description |
|-------|-------------|
| `GET /` | Page d'accueil avec présentation de l'application |
| `GET /tracker` | Liste de tous les artistes avec recherche et filtres |
| `GET /artist/{id}` | Page détaillée d'un artiste (membres, concerts, dates, lieux) |
| `GET /locations` | Carte interactive des lieux de concerts |
| `GET /compare` | Page de comparaison de deux artistes |
| `GET /favorites` | Liste des artistes favoris de l'utilisateur |
| `GET /favorite?id={id}&back={page}` | Toggle favori (ajouter/retirer) |

## Fonctionnalités implémentées

### Fonctionnalités obligatoires

#### 1. Page d'accueil
* Page d'accueil élégante avec présentation de l'application
* Bouton "Entrer dans le site" menant vers la liste des artistes
* Navigation claire et intuitive

#### 2. Liste des artistes
* Affichage de tous les artistes en grille responsive
* Informations affichées par carte : image, nom, année de création, nombre de membres
* Lien cliquable vers la page détaillée de chaque artiste

#### 3. Page de détails d'un artiste
* Affichage complet : image, nom, année de création, premier album
* Liste de tous les membres du groupe
* Historique détaillé des concerts avec dates et lieux
* Navigation fluide avec bouton retour

#### 4. Recherche
* Barre de recherche fonctionnelle sur `/tracker`
* Recherche par nom d'artiste avec préfixe (ex: "Que" trouve "Queen")
* Traitement 100% côté Go via requête HTTP GET
* Affichage dynamique des résultats

#### 5. Filtres
* Filtre par intervalle : année de création (min/max)
* Filtre par sélection multiple : nombre de membres (1 à 8+)
* Combinaison possible de tous les filtres
* Boutons "Apply" et "Reset" pour gérer les filtres
* Tout traité côté serveur

#### 6. Événement interactif
* Clic sur un lieu dans la liste des concerts → scroll automatique vers ce lieu sur la carte
* Clic sur un artiste dans la carte → redirection vers sa page détaillée

#### 7. Gestion des erreurs
* Page 404 personnalisée pour les routes inexistantes
* Gestion des IDs d'artistes invalides
* Pas de crash serveur - toutes les erreurs sont gérées proprement
* Messages d'erreur clairs pour l'utilisateur

---

### **TOUS LES BONUS ONT ÉTÉ RÉALISÉS**

#### Bonus 1 : Visualisations (Carte interactive)
  **IMPLÉMENTÉ**

* Carte mondiale interactive affichant tous les lieux de concerts
* 186 lieux uniques répartis sur 5 continents
* Projection équirectangulaire pour convertir lat/lon en coordonnées
* Marqueurs SVG cliquables avec animation au survol
* Système de zoom (1x à 4x) avec boutons +/- et reset
* Navigation par drag & drop sur la carte zoomée
* Regroupement automatique par continents
* Liste détaillée sous la carte organisée par continent
* Compteurs de concerts par continent et par lieu
* Affichage des artistes par lieu au clic

**Page** : `/locations`

#### Bonus 2 : Filtres
  **IMPLÉMENTÉ**

* Filtre par intervalle d'années (min/max)
* Filtre par nombre de membres (checkboxes 1-8+)
* Combinaison de plusieurs filtres simultanément
* Bouton "Apply" pour appliquer les filtres
* Bouton "Reset" pour réinitialiser
* Tout le traitement est fait côté serveur Go
* Aucune utilisation de JavaScript pour les filtres

**Page** : `/tracker` (menu "Filtres")

#### Bonus 3 : Barre de recherche
   **IMPLÉMENTÉE**

* Recherche par nom d'artiste
* Recherche par préfixe (commence par...)
* Traitement côté serveur uniquement
* Requête HTTP GET standard
* Résultats affichés dynamiquement
* Compatible avec les filtres

**Page** : `/tracker` (barre de recherche en haut)

#### Bonus 4 : Favoris
  **IMPLÉMENTÉ**

* Système complet de gestion des favoris
* Ajout/suppression d'artistes favoris via icône cœur
* Page dédiée `/favorites` listant tous les favoris
* Persistance des favoris via cookies
* Icônes interactives changeant d'état (cœur rempli/vide)
* Fonction de retour intelligente (retour à la page d'origine après toggle)
* Compteur visible du nombre de favoris

**Pages** : 
- `/tracker` (icônes sur chaque carte)
- `/favorites` (liste complète)

#### Bonus 5 : Comparaison
   **IMPLÉMENTÉE**

* Page `/compare` permettant de comparer deux artistes côte à côte
* Menus déroulants pour sélectionner les artistes
* Badge "VS" stylisé entre les deux artistes
* Comparaison détaillée affichant :
  - Photos des artistes
  - Année de création
  - Premier album
  - Nombre de membres
  - Liste complète des membres
  - Nombre de concerts
  - Lieux de concerts (avec scroll)
* Affichage en grille responsive (2 colonnes)
* Design cohérent avec le reste du site

**Page** : `/compare`

---

## Structure du projet

```
groupie-tracker/
├── main.go                # Point d'entrée et configuration des routes
├── handlers.go            # Gestionnaires HTTP pour toutes les pages
├── types.go               # Structures de données (Artist, Location, etc.)
├── utils.go               # Fonctions utilitaires (formatage, cookies)
├── api.go                 # Appels à l'API Groupie Trackers
├── geo.go                 # Coordonnées géographiques et continents
├── go.mod                 # Dépendances Go
├── README.md              # Documentation
├── templates/
│   ├── welcome.html       # Page d'accueil
│   ├── index.html         # Liste des artistes
│   ├── artist.html        # Détails d'un artiste
│   ├── locations.html     # Carte interactive (Bonus 1)
│   ├── compare.html       # Comparaison d'artistes (Bonus 5)
│   └── favorites.html     # Artistes favoris (Bonus 4)
└── static/
    ├── css/
    │   ├── style.css      # Styles globaux
    │   ├── locations.css  # Styles de la carte
    │   └── compare.css    # Styles de comparaison
    └── img/
        └── light_logo.png # Logo de l'application
```

## Architecture technique

### Backend Go
* **Package net/http** pour le serveur et les routes
* **Package html/template** pour le rendu des pages dynamiques
* **Package encoding/json** pour parser les données API
* Organisation modulaire en 6 fichiers :
  - `main.go` : Routes et configuration (30 lignes)
  - `handlers.go` : Logique métier des pages (400 lignes)
  - `types.go` : Structures de données (90 lignes)
  - `utils.go` : Fonctions utilitaires (110 lignes)
  - `api.go` : Communication avec l'API externe (30 lignes)
  - `geo.go` : Traitement géographique (120 lignes)

### API externe
* Base URL : `https://groupietrackers.herokuapp.com/api`
* Endpoints consommés :
  - `/artists` - Liste complète des 52 artistes
  - `/locations/{id}` - Lieux de concerts d'un artiste
  - `/dates/{id}` - Dates de concerts d'un artiste
  - `/relation/{id}` - Relations lieux-dates

### Frontend
* HTML5 avec templates Go
* CSS3 personnalisé avec palette de couleurs élégante (vert forêt/beige)
* Police Playfair Display pour un rendu professionnel
* Design responsive adaptatif
* SVG pour les icônes et la carte interactive
* JavaScript minimal (uniquement pour zoom/drag de la carte)

## Fonctionnalités techniques avancées

### Système de cookies
* Stockage des favoris dans un cookie nommé "favorites"
* Format : IDs séparés par des virgules (ex: "1,3,5")
* Persistance entre les sessions
* Lecture/écriture sécurisée

### Traitement géographique
* Conversion de 50+ noms de pays en coordonnées latitude/longitude
* Projection équirectangulaire pour affichage sur carte 2D
* Regroupement automatique par continents (Amérique du Nord/Sud, Europe, Asie, Océanie, Autre)
* Calcul des coordonnées X/Y pour positionnement précis des marqueurs

### Templates Go
* Utilisation du package `html/template` avec fonctions personnalisées :
  - `formatPlace` : Formate les noms de lieux (ville, pays)
  - `formatDate` : Nettoie les dates (supprime les astérisques)
* Rendu dynamique basé sur les données de l'API
* Protection contre les injections XSS

### Gestion d'état
* Pas de base de données nécessaire
* Toutes les données sont récupérées de l'API à la volée
* Cache en mémoire pendant la durée de vie d'une requête
* État des favoris persisté via cookies côté client

## Tests effectués

### Tests fonctionnels
```bash
# Test recherche
# → Taper "Queen" dans la barre de recherche
# → Vérifier que seul Queen apparaît

# Test filtres
# → Année min: 1970, max: 1980
# → Nombre de membres: 4
# → Cliquer Apply
# → Vérifier que seuls les groupes de 4 membres créés entre 1970-1980 apparaissent

# Test carte
# → Aller sur /locations
# → Zoomer avec les boutons +/-
# → Cliquer sur un marqueur
# → Vérifier le scroll vers la liste

# Test comparaison
# → Aller sur /compare
# → Sélectionner "Queen" et "Pink Floyd"
# → Cliquer "COMPARER"
# → Vérifier l'affichage côte à côte

# Test favoris
# → Cliquer sur l'icône cœur de Queen
# → Aller sur /favorites
# → Vérifier que Queen apparaît
# → Cliquer à nouveau sur le cœur
# → Vérifier la suppression
```

### Tests de robustesse
* URL invalide → 404
* ID d'artiste inexistant → 404
* Filtres avec valeurs extrêmes → Gestion propre
* Cookies désactivés → Favoris non sauvegardés mais pas de crash
* API indisponible → Message d'erreur approprié

## Design

L'application utilise une palette de couleurs inspirée de la nature :
* **Vert forêt** (#3d4e35) pour la sidebar et les cartes
* **Beige** (#e8dcc8) pour les textes
* **Vert clair** (#8fbc8f) pour les accents et hover
* **Police élégante** Playfair Display pour un look professionnel et classique

Effet visuel :
* Animations douces au survol
* Transitions fluides entre les pages
* Responsive design (desktop, tablette, mobile)
* Icônes SVG vectorielles

## Récapitulatif des bonus

| Bonus | Statut | Page/Fonctionnalité |
|-------|--------|---------------------|
| Visualisations (Carte) |  Implémenté | `/locations` |
| Filtres |  Implémenté | `/tracker` (menu Filtres) |
| Barre de recherche |  Implémenté | `/tracker` (barre en haut) |
| Favoris |  Implémenté | `/favorites` + icônes partout |
| Comparaison |  Implémenté | `/compare` |

**TOUS LES BONUS DEMANDÉS ONT ÉTÉ RÉALISÉS ET SONT FONCTIONNELS.**

---

Bon voyage musical !
