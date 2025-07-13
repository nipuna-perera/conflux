// @ts-check

import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import svelte from "eslint-plugin-svelte";
import prettier from "eslint-config-prettier";
import globals from "globals";
import unusedImports from "eslint-plugin-unused-imports";

export default tseslint.config(
  // Global ignores
  {
    ignores: [
      "build/",
      ".svelte-kit/",
      "dist/",
      "node_modules/",
      ".dagger/",
      "**/*.min.js",
      "**/vendor/**",
      "coverage/",
      ".env*",
      "*.config.js",
      "*.config.ts"
    ],
  },
  
  // Base JavaScript/TypeScript configuration
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  
  // Base configuration for all files
  {
    files: ["**/*.{js,ts,svelte}"],
    plugins: {
      "unused-imports": unusedImports,
    },
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
      parser: tseslint.parser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: "module",
        extraFileExtensions: [".svelte"],
      },
    },
    rules: {
      // Import/Export management
      "unused-imports/no-unused-imports": "error",
      "unused-imports/no-unused-vars": [
        "error",
        {
          "vars": "all",
          "varsIgnorePattern": "^_",
          "args": "after-used",
          "argsIgnorePattern": "^_",
          "caughtErrors": "all",
          "caughtErrorsIgnorePattern": "^_",
          "destructuredArrayIgnorePattern": "^_"
        }
      ],
      
      // TypeScript specific rules
      "@typescript-eslint/no-unused-vars": "off", // Handled by unused-imports
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/no-empty-object-type": "warn",
      "@typescript-eslint/no-require-imports": "off",
      
      // General JavaScript rules
      "no-console": ["warn", { "allow": ["warn", "error"] }],
      "no-debugger": "error",
      "no-unreachable": "error",
      "no-unused-expressions": "error",
      "prefer-const": "error",
      "prefer-template": "warn",
      
      // Disable conflicting rules
      "no-undef": "off", // TypeScript handles this
    },
  },
  
  // Svelte-specific configuration
  ...svelte.configs["flat/recommended"],
  {
    files: ["**/*.svelte"],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
        extraFileExtensions: [".svelte"],
      },
    },
    rules: {
      // Svelte specific rules  
      "svelte/no-at-debug-tags": "warn",
      "svelte/no-at-html-tags": "warn",
      "svelte/require-each-key": "error",
      "svelte/valid-compile": "error",
      "svelte/no-unused-svelte-ignore": "warn",
      
      // Allow some flexibility for Svelte components
      "unused-imports/no-unused-vars": [
        "error",
        {
          "vars": "all",
          "varsIgnorePattern": "^(_|\\$\\$)",
          "args": "after-used",
          "argsIgnorePattern": "^_",
          "caughtErrors": "all",
          "caughtErrorsIgnorePattern": "^_"
        }
      ],
    },
  },
  
  // SvelteKit-specific files
  {
    files: [
      "**/app.d.ts",
      "**/hooks.client.ts",
      "**/hooks.server.ts",
      "**/+*.ts",
      "**/+*.js"
    ],
    rules: {
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
          caughtErrorsIgnorePattern: "^_",
          destructuredArrayIgnorePattern: "^_"
        }
      ],
    },
  },
  
  // Configuration files
  {
    files: ["vite.config.{js,ts}", "svelte.config.js", "tailwind.config.{js,ts}"],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
  },
  
  // Prettier configuration (must be last)
  prettier,
  ...svelte.configs["flat/prettier"]
);
