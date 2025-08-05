# Code Style and Conventions

## TypeScript Configuration
- **Target**: ESNext
- **Strict mode**: Enabled
- **JSX**: react-jsx transform
- **Module**: ESNext with Node resolution
- **No emit**: true (build handled by Vite)

## ESLint Configuration
- Extends: eslint:recommended, @typescript-eslint/recommended, react-hooks/recommended
- Parser: @typescript-eslint/parser
- Plugins: react-refresh
- Environment: browser, es2020

## File Structure Conventions
- **Components**: `/src/components/[component-name]/index.tsx`
- **Pages**: `/src/pages/[resource]/[action].tsx` (list.tsx, create.tsx, edit.tsx, show.tsx)
- **Contexts**: `/src/contexts/[context-name]/index.tsx`
- **Providers**: Root level files (dataProvider.ts, authProvider.ts)

## Naming Conventions
- Files: kebab-case for directories, camelCase for TypeScript files
- Components: PascalCase
- Functions: camelCase
- Constants: camelCase (following the existing dataProvider pattern)

## Import/Export Patterns
- Default exports for main components
- Named exports collected in index.ts files for cleaner imports
- ES modules (type: "module" in package.json)