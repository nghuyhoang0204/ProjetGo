# Guide de démarrage - Transpilateur TypeScript → JavaScript

## 🚀 Introduction

Ce projet est un transpilateur TypeScript vers JavaScript écrit en Go pur, sans aucune dépendance externe. Il permet de convertir du code TypeScript en code JavaScript équivalent.

## 💻 Comment exécuter le projet

### Option 1 : Interface Web (recommandé)

```bash
# Compiler le projet
go build -o ProjetGO.exe

# Exécuter l'application
./ProjetGO.exe
```

L'interface web sera disponible à l'adresse [http://localhost:8080](http://localhost:8080)

### Option 2 : Mode Console

```bash
# Avec un exemple par défaut
./ProjetGO.exe -console

# Avec un fichier TypeScript spécifique
./ProjetGO.exe -console -file chemin/vers/votre/fichier.ts
```

### Options disponibles

```
-port string     Port pour le serveur web (défaut "8080")
-host string     Host pour le serveur web (défaut "localhost")
-console         Exécution en mode console (défaut false)
-file string     Fichier d'entrée à transpiler
-verbose         Affichage détaillé (défaut false)
```

## 🧪 Tests

Pour exécuter les tests de validation du transpilateur :

```bash
# Test des paramètres par défaut
go run test_default_params.go

# Autres tests disponibles
go run test_function.go
```

## 📋 Exemples d'utilisation

### Exemple 1 : Fonction avec paramètre par défaut

```typescript
// TypeScript
function incrementer(nombre: number, pas: number = 1): number {
  return nombre + pas;
}

// JavaScript généré
function incrementer(nombre, pas = 1) {
  return nombre + pas;
}
```

### Exemple 2 : Classe avec méthodes

```typescript
// TypeScript
class Calculator {
  private value: number = 0;
  
  add(x: number): void {
    this.value += x;
  }
  
  getResult(): number {
    return this.value;
  }
}

// JavaScript généré
class Calculator {
  constructor() {
    this.value = 0;
  }
  
  add(x) {
    this.value += x;
  }
  
  getResult() {
    return this.value;
  }
}
```

### Exemple 3 : Manipulation de dates

```typescript
// TypeScript
function formaterDate(date: Date): string {
  const jour = date.getDate().toString().padStart(2, '0');
  const mois = (date.getMonth() + 1).toString().padStart(2, '0');
  const annee = date.getFullYear();
  
  return jour + '/' + mois + '/' + annee;
}

// JavaScript généré
function formaterDate(date) {
  const jour = date.getDate().toString().padStart(2, '0');
  const mois = (date.getMonth() + 1).toString().padStart(2, '0');
  const annee = date.getFullYear();
  
  return jour + '/' + mois + '/' + annee;
}
```

### Exemple 4 : Classe avec manipulation temporelle

```typescript
// TypeScript
class GestionnaireDates {
  private dates: Date[];
  
  constructor(dateInitiale: Date = new Date()) {
    this.dates = [dateInitiale];
  }
  
  ajouterJours(date: Date, jours: number = 1): Date {
    const nouvelleDate = new Date(date);
    nouvelleDate.setDate(date.getDate() + jours);
    return nouvelleDate;
  }
}

// JavaScript généré
class GestionnaireDates {
  constructor(dateInitiale = new Date()) {
    this.dates = [dateInitiale];
  }
  
  ajouterJours(date, jours = 1) {
    const nouvelleDate = new Date(date);
    nouvelleDate.setDate(date.getDate() + jours);
    return nouvelleDate;
  }
}
```

## 📝 Structure du projet

```
go.mod              # Fichier de configuration Go
main.go             # Point d'entrée principal
web.go              # Interface web
ast/                # Module de représentation syntaxique abstraite
  ast.go            # Définition des structures AST
lexer/              # Module d'analyse lexicale
  lexer.go          # Tokenisation du code source
parser/             # Module d'analyse syntaxique
  parser.go         # Construction de l'AST
generator/          # Module de génération de code
  generator.go      # Générateur de code JavaScript
  utils.go          # Fonctions utilitaires pour la génération
```

## ⚠️ Limitations actuelles

- Les fonctionnalités TypeScript avancées comme les génériques complexes peuvent ne pas être correctement transpilées
- Certains cas particuliers d'héritage de classe peuvent nécessiter des améliorations
