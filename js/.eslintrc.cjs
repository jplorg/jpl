module.exports = {
  env: { browser: true, node: true, es2021: true },
  overrides: [
    {
      env: { node: true },
      files: ['.eslintrc.{js,cjs}'],
      parserOptions: { sourceType: 'script' },
    },
  ],
  extends: ['eslint:recommended', 'plugin:prettier/recommended'],
  parserOptions: { ecmaVersion: 'latest', sourceType: 'module' },
  plugins: ['unused-imports', 'import'],
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
};
