# Task Completion Checklist

## When a coding task is completed, run these commands:

### 1. Type Checking
```bash
npx tsc --noEmit
```
- Ensures TypeScript compilation without errors
- Required since strict mode is enabled

### 2. Linting
```bash
npx eslint src --ext .ts,.tsx
```
- Checks code style and potential issues
- Follows the configured ESLint rules

### 3. Build Verification
```bash
yarn build
```
- Runs TypeScript compilation and Refine build
- Ensures production build works correctly

### 4. Development Server Test (Optional)
```bash
yarn dev
```
- Start dev server to manually verify changes work as expected
- Particularly important for UI changes

## Notes
- The project uses Yarn as package manager
- No specific test scripts are configured in package.json
- Vite handles the build process with hot reload in development
- Refine dev server provides additional development features