module.exports = {
    "env": {
        "browser": true,
        "es2021": true
    },
    "extends": [
        "eslint:recommended",
        "plugin:vue/vue3-essential",
        "plugin:@typescript-eslint/recommended",
        "prettier",
        "plugin:vue/base",
        "plugin:vue/vue3-recommended"
    ],
    "overrides": [
    ],
    "parser": "vue-eslint-parser",

    "parserOptions": {
        "parser": "@typescript-eslint/parser",
        "ecmaVersion": "latest",
        "sourceType": "module"
    },
    "plugins": [
        "vue",
        "@typescript-eslint"
    ],
    "rules": {
        "@typescript-eslint/ban-ts-comment":"off",
        "vue/multi-word-component-names":"off",
        "@typescript-eslint/no-empty-function":"off",
        "@typescript-eslint/no-explicit-any":"off"
    }
}
