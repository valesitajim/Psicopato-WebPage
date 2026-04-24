# Psicopato-web Page

## Overview

Psicopato-web Page is a modern, interactive, and highly optimized e-commerce web application built with Python and Flask, tailored for clothing and duck-themed merchandise. This storefront provides a fully responsive layout, combining a clean user interface with advanced web performance optimizations. It is designed to offer a seamless shopping experience across all devices, utilizing scalable SCSS architecture (BEM methodology) and Vanilla JavaScript for smooth interactivity.

## Features

* Fully responsive layout adapted for mobile, tablet, and desktop screens
* Interactive UI elements including custom dropdown menus (Profile, Favorites, Cart) built with Vanilla JavaScript
* High-performance image delivery using lightweight AVIF formats and `<picture>` tags for legacy browser support
* Asynchronous loading of Google Fonts to eliminate render-blocking and improve Lighthouse scores
* Scalable and maintainable stylesheet architecture using SCSS

## Dependencies 

* A modern web browser (Chrome, Firefox, Safari, Edge)
* Live Sass Compiler extension in VS code
* Golang 

## New Features and Functionalities

In this second iteration of the project, the static interface has been migrated to a dynamic backend server using **Go (Golang)**, applying an **MVC (Model-View-Controller)** architecture.

### Security and Authentication
* **Password Encryption:** Integration of the official library `golang.org/x/crypto/bcrypt` to generate secure hashes (cost 12). Passwords are never stored or logged in plain text.

* **Data Leak Prevention:** Login error messages are generic ("Incorrect Credentials") to prevent user enumeration.

* **Secure Logging:** A categorized console event logging system (`INFO`, `WARN`, `ERROR`) monitors access and failures without exposing sensitive information or secrets.

### Backend and Routing (Go)
* **Separation of Responsibilities (MVC):** * `models`: Definition of data structures (`User`).

* `database`: Management of simulated persistence using `.jsonl` files.

* `handlers`: Controllers that process `GET` and `POST` methods.


### UX/UI and Frontend Improvements
* **Notification System (Toasts):** Non-intrusive pop-up alerts based on URL parameters to report successful logins or errors.

* **Accessibility and Micro-interactions:** Implementation of an interactive button (with pure SVG icons) to show or hide the password in the registration form.

---

## Installation and Execution Instructions

1. **Clone the repository and navigate to the root folder.**
2. **Install dependencies:**

Since the cryptography library has been added, it is necessary to synchronize the Go modules by running:
``bash
go mod tidy

3. **In terminal run:** 
``bash
go run ./cmd/web/main.go

## License 
This project is licensed under the MIT License. Feel free to use, modify, and distribute it
