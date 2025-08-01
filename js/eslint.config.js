import js from '@eslint/js';
import importPlugin from 'eslint-plugin-import';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import unusedImportsPlugin from 'eslint-plugin-unused-imports';
import globals from 'globals';

export default [
  js.configs.recommended,
  eslintPluginPrettierRecommended,

  { ignores: ['lib/*'] },
  { files: ['src/**/*.{js,jsx}'] },
  {
    linterOptions: {
      reportUnusedDisableDirectives: 'warn',
    },
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2022,
      },
    },
    plugins: {
      'unused-imports': unusedImportsPlugin,
      import: importPlugin,
    },
    rules: {
      'no-unused-vars': ['error', { ignoreRestSiblings: true }],
      // 'no-use-before-define': 'error',
      'arrow-body-style': 'error',
      'unused-imports/no-unused-imports': 'error',
      'import/no-unresolved': 'error',
      'import/no-duplicates': 'error',
      'import/order': [
        'error',
        {
          groups: ['builtin', 'external', 'internal', 'parent', ['sibling', 'index']],
          alphabetize: { order: 'asc', orderImportKind: 'asc' },
          'newlines-between': 'never',
        },
      ],
      'import/newline-after-import': 'error',
    },
  },
];
