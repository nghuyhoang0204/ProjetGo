package main

// 🚀 TRANSPILEUR TYPESCRIPT → JAVASCRIPT 
// =====================================
// 
// ✨ Projet Go PURE - ZERO dépendance externe
// 
// Ce transpileur est écrit en Go pur, utilisant uniquement la bibliothèque standard :
// - net/http : Serveur web intégré
// - encoding/json : Parsing/génération JSON
// - strings : Manipulation de chaînes
// - html/template : Templates HTML
// - fmt, os, io : Opérations de base
//
// 🎯 Objectif : Transpiler du code TypeScript en JavaScript idiomatique
// 🏗️ Architecture modulaire : lexer → parser → AST → generator
// 🌐 Interface web intégrée pour tester interactivement
//
// 📁 Structure :
// - ast/     : Définition de l'Abstract Syntax Tree  
// - lexer/   : Analyseur lexical (tokens)
// - parser/  : Analyseur syntaxique (AST)
// - generator/ : Générateur de code JavaScript
// - main.go  : Point d'entrée et serveur web
//
// 🔧 Ce fichier sert maintenant de documentation
// Toutes les fonctionnalités sont intégrées dans main.go
