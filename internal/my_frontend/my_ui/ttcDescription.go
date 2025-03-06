package myapp

import (
	"fmt"
)

var deDescription string = fmt.Sprintln(`Beschreibung der Anwendung:
TT Companion ist eine Anwendung zur Verwaltung von Tischtennis-Clubs. Der Benutzer kann Clubs, Teams, Spieler erstellen und Beziehungen zwischen ihnen hinzufügen.
Alle neuen Änderungen werden auf der Client-Seite vorgenommen und müssen vor dem Beenden gespeichert werden, entweder über das Hauptmenü oder auf der Hauptseite.
Der Benutzer kann auch seine Sprach- und Designpräferenzen in den Optionen anpassen. Vergessen Sie nicht, Ihre Arbeit vor dem Verlassen der App zu speichern!
Ziel der Anwendung:
Der Zweck dieser Anwendung ist es, meine Fähigkeiten in der Softwareentwicklung zu stärken:
    • Backend-Sprache: Golang 
    • Frontend-Sprache: Golang (Framework Fyne) 
    • Versionsverwaltung: Git 
    • Verteilung: Über ein Docker-Image, das eine komprimierte Version der kompilierten Software in WebAssembly (WASM) enthält. Die Komprimierung erfolgt mit Brotli. 
    • Speicherung: Github und Dockerhub. https://github.com/Whadislov/TTCompanion 
    • Datenbank: PostgreSQL, bereitgestellt von Neon.tech. Gespeicherte Passwörter sind mit der Funktion bcrypt verschlüsselt. 
    • Automatisierung: Hauptsächlich GitHub Actions für Bereitstellungen, etwas Jenkins für Tests und Kompilierungen. 
    • Bereitstellung: 
        ◦ Die Staging-Umgebung läuft auf Render.com. https://ttcompanion.onrender.com 
        ◦ Die Produktionsumgebung läuft auf Cloud Run. https://ttcompanion-prod-912172190800.europe-west9.run.app`)

var frDescription string = fmt.Sprintln(`Description de l’application :
TT Companion est une application de gestion de clubs de tennis de table. L’utilisateur peut créer des clubs, des équipes, des joueurs et établir des relations entre eux.
Toutes les modifications sont effectuées côté client et doivent être enregistrées avant de quitter l’application, soit via le menu principal, soit sur la page d’accueil.
L’utilisateur peut également ajuster ses préférences de langue et de thème dans les options. N’oubliez pas de sauvegarder votre travail avant de quitter l’application !
Objectif de l’application :
Le but de cette application est de renforcer mes compétences en développement logiciel :
    • Langage backend : Golang 
    • Langage frontend : Golang (Framework Fyne) 
    • Gestion de versions : Git 
    • Distribution : Via une image Docker contenant une version compressée du logiciel compilé en WebAssembly (WASM). La compression est réalisée avec Brotli. 
    • Stockage : Github et Dockerhub. https://github.com/Whadislov/TTCompanion 
    • Base de données : PostgreSQL, fourni par Neon.tech. Les mots de passe stockés sont chiffrés avec la fonction bcrypt. 
    • Automatisation : Principalement GitHub Actions pour les déploiements, un peu de Jenkins pour les tests et les compilations. 
    • Déploiement : 
        ◦ L’environnement de staging fonctionne sur Render.com. https://ttcompanion.onrender.com 
        ◦ L’environnement de production fonctionne sur Cloud Run. https://ttcompanion-prod-912172190800.europe-west9.run.app`)

var enDescription string = fmt.Sprintln(`Description of the Application
TT Companion is a table tennis club management application. The user can create clubs, teams, players, and establish relationships between them.
All modifications are performed on the client side and must be saved before exiting, either via the main menu or on the main page.
Users can also customize their language and theme preferences in the options. Don’t forget to save your work before quitting the app!
Goal of the Application
The purpose of this application is to enhance my software development skills:
    • Backend language: Golang 
    • Frontend language: Golang (Framework Fyne) 
    • Version management: Git 
    • Distribution: Via a Docker image containing a compressed version of the compiled software in WebAssembly (WASM). Compression is handled with Brotli. 
    • Storage: Github and Dockerhub. https://github.com/Whadislov/TTCompanion 
    • Database: PostgreSQL, provided by Neon.tech. Stored passwords are encrypted using the bcrypt function. 
    • Automation: Primarily GitHub Actions for deployments, with some Jenkins for tests and compilations. 
    • Deployment: 
        ◦ The staging environment runs on Render.com. https://ttcompanion.onrender.com 
        ◦ The production environment runs on Cloud Run. https://ttcompanion-prod-912172190800.europe-west9.run.app`)
