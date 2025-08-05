# Codebase Structure

## Root Structure
```
/
├── public/           # Static assets
├── src/             # Source code
├── .serena/         # Serena configuration
├── .claude/         # Claude configuration
├── package.json     # Dependencies and scripts
├── tsconfig.json    # TypeScript configuration
├── vite.config.ts   # Vite build configuration
├── .eslintrc.cjs    # ESLint configuration
└── README.MD        # Project documentation
```

## Source Code Structure (`/src`)
```
src/
├── App.tsx                 # Main application component
├── index.tsx              # Application entry point
├── dataProvider.ts        # API data provider (REST)
├── authProvider.ts        # Authentication provider
├── vite-env.d.ts         # Vite TypeScript definitions
├── components/           # Reusable components
│   ├── header/          # Header component
│   └── index.ts         # Component exports
├── contexts/            # React contexts
│   └── color-mode/      # Theme/color mode context
├── pages/               # Page components (CRUD operations)
│   ├── blog-posts/      # Blog post management
│   │   ├── list.tsx     # List view
│   │   ├── create.tsx   # Create form
│   │   ├── edit.tsx     # Edit form
│   │   ├── show.tsx     # Detail view
│   │   └── index.ts     # Page exports
│   ├── categories/      # Category management
│   │   └── [same structure as blog-posts/]
│   ├── login/           # Login page
│   ├── register/        # Registration page
│   └── forgotPassword/  # Password recovery
```

## Key Features
- **Blog Posts**: Full CRUD operations for blog post management
- **Categories**: Category management (currently with mock data)
- **Authentication**: Login, register, forgot password flows
- **Responsive UI**: Material-UI components with color mode support
- **Data Grid**: Advanced table functionality for list views