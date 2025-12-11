# Power4-Web

Une application web interactive du jeu **Puissance 4** construite avec Go et servie via HTTP.

## 🎮 Fonctionnalités

- Jeu Puissance 4 intégralement fonctionnel
- Mode sombre/clair avec sauvegarde des préférences
- Système de pouvoirs spéciaux
- Serveur léger en Go

## 🚀 Installation et démarrage

1. Clonez le repository :
```bash
git clone https://github.com/CyrilMaz/power4-web.git
cd power4-web
```

2. Lancez le serveur :
```bash
go run main.go
```

3. Accédez à l'application :
```
http://localhost:8080
```

## 📁 Structure du projet

```
power4-web/
├── main.go              # Point d'entrée de l'application
├── game/                # Logique du jeu Puissance 4
├── handlers/            # Gestionnaires HTTP
├── theme/               # Gestion du thème (clair/sombre)
├── static/              # Fichiers statiques (CSS, JS)
├── templates/           # Fichiers HTML
└── go.mod              # Dépendances Go
```

## 🎮 Comment jouer

1. Sélectionnez une colonne pour y placer votre pion
2. Gagnez en alignant 4 pions de votre couleur
3. Utilisez les pouvoirs spéciaux pour des coups stratégiques
4. Réinitialisez la partie avec le bouton "Reset"

## 🌓 Thème

Basculez entre le mode clair et sombre avec le bouton de bascule en haut de la page. Votre préférence est automatiquement sauvegardée.

## 👤 Auteur

Robbe Matthias
Cyril Mazauric
Nathan Gueroult
