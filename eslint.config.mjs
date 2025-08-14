import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import playwright from 'eslint-plugin-playwright';

export default [
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['**/*.{js,ts}'],
    ignores: [
      'node_modules/',
      'dist/',
      'build/',
      'out/',
      'test-results/',
      'playwright-report/',
      'coverage/',
      'books_api/',
      '*.min.js',
      '*.bundle.js',
    ],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        browser: true,
        node: true,
        es2021: true,
      },
    },
    rules: {
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/no-explicit-any': 'warn',
      'no-console': 'warn',
    },
  },
  {
    files: ['tests/**/*.ts', '**/*.spec.ts', '**/*.test.ts'],
    plugins: {
      playwright: playwright,
    },
    rules: {
      ...playwright.configs.recommended.rules,
      'no-console': 'off',
      'playwright/expect-expect': 'error',
      'playwright/max-nested-describe': ['error', { max: 3 }],
      'playwright/no-page-pause': 'error',
      'playwright/no-element-handle': 'error',
      'playwright/no-eval': 'error',
      'playwright/no-focused-test': 'error',
      'playwright/no-force-option': 'warn',
      'playwright/no-skipped-test': 'warn',
      'playwright/no-useless-await': 'error',
      'playwright/prefer-to-have-length': 'error',
      'playwright/prefer-web-first-assertions': 'error',
    },
  },
];
