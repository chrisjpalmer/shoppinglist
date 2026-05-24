# Architecture

## Overview

- DB: SQLite
- Backend: Go
    - Proto
    - Gotempl
- Frontend: Svelte
    - Tailwindcss

- Deployment:
    - k3s on a raspberry pi
    - Railway for the preview site

- CI/CD: Github Actions + Dagger

## Introduction

The architecture chosen for this project is deliberately simple to keep maintenance low.
Performance and scaling were not considered, since I am the only user.

## Database

For this project **SQLite** was used. The advantages using SQLite over something like
Postgres / MySQL, is that it is simply one file on the disk. You don't need to deploy
a dedicated database anywhere because the SQLite library handles everything - all it needs
is the file.

## Backend

The backend is written in Go. The backend has an API server, which services the frontend *Meal Planning Site*, 
and a *Shopping Site* server.

The API uses protobuf over connectrpc which gives the benefits of protobuf without having to use the binary protocol. 
Instead the backend / frontend communicate via JSON payloads. This makes debugging much easier.

The Shopping Site server is backend by gotempl and htmx. I chose to use something different just for fun.
Gotempl is an html templating language that generates go and htmx is a very lightweight library for
the browser that can swap the html out for parts of your DOM. The appeal for using something like htmx
is that it eliminates a frontend framework such as react of svelte. It also eliminates the need to use
protocol buffers. This is a very novel concept in the current software landscape where we are continually
reaching for frontend frameworks - I quite like htmx so far and want to see how far it goes.

## Frontend - Meal Planning Site

The frontend, which only hosts the *Meal Planning Site* is written in Svelte Kit. I haven't really had a lot of
exposure to other frameworks so the decision to use Svelte was purely out of recommendations from other people.
The things I liked about Svelte is its simplicity.

## CI/CD

I used a combination of Github Actions and Dagger for CI/CD. Github Actions is used just to "kick things off"
whilst Dagger is used to do all the heavy lifting. The thing I like about Dagger is that you can write CI/CD
pipelines that run both locally and in Github Actions. This makes it easy to test them.