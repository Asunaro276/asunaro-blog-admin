# Asunaro CMS - Project Overview

## Purpose
This is a frontend admin panel for the Asunaro blog system built with Refine framework. It's designed as an internal tool for managing blog posts and categories.

## Tech Stack
- **Framework**: React 18 with TypeScript
- **UI Framework**: Refine.dev with Material-UI (MUI)
- **Build Tool**: Vite
- **Routing**: React Router v7
- **State Management**: React Hook Form
- **Data Grid**: MUI X Data Grid
- **Package Manager**: Yarn (based on yarn.lock presence)

## Key Dependencies
- **@refinedev/core**: Main Refine framework
- **@refinedev/mui**: Material-UI integration for Refine
- **@refinedev/simple-rest**: REST data provider
- **@refinedev/react-hook-form**: Form management integration
- **@mui/material**: Material-UI components
- **@mui/x-data-grid**: Advanced data grid component

## Architecture
The application follows Refine's architecture pattern with:
- Data providers for API communication (currently with mock implementations)
- Auth providers for authentication
- Pages for CRUD operations on blog posts and categories
- Color mode context for theming