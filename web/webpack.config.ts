// eslint-disable-next-line @typescript-eslint/no-triple-slash-reference, spaced-comment
/// <reference path="../shared/src/types/terser-webpack-plugin/index.d.ts" />

import MiniCssExtractPlugin from 'mini-css-extract-plugin'
import MonacoWebpackPlugin from 'monaco-editor-webpack-plugin'
import OptimizeCssAssetsPlugin from 'optimize-css-assets-webpack-plugin'
import * as path from 'path'
import TerserPlugin from 'terser-webpack-plugin'
import * as webpack from 'webpack'

const mode = process.env.NODE_ENV === 'production' ? 'production' : 'development'
console.error('Using mode', mode)

const devtool = mode === 'production' ? 'source-map' : 'cheap-module-eval-source-map'

const rootDir = path.resolve(__dirname, '..')
const nodeModulesPath = path.resolve(__dirname, '..', 'node_modules')
const monacoEditorPaths = [path.resolve(nodeModulesPath, 'monaco-editor')]

const isEnterpriseBuild = !!process.env.ENTERPRISE
const enterpriseDir = path.resolve(__dirname, 'src', 'enterprise')

const config: webpack.Configuration = {
    context: __dirname, // needed when running `gulp webpackDevServer` from the root dir
    mode,
    optimization: {
        minimize: mode === 'production',
        minimizer: [
            new TerserPlugin({
                sourceMap: true,
                terserOptions: {
                    compress: {
                        // // Don't inline functions, which causes name collisions with uglify-es:
                        // https://github.com/mishoo/UglifyJS2/issues/2842
                        inline: 1,
                    },
                },
            }),
        ],
        namedModules: false,

        ...(mode === 'development'
            ? {
                  removeAvailableModules: false,
                  removeEmptyChunks: false,
                  splitChunks: false,
              }
            : {}),
    },
    entry: {
        // Enterprise vs. OSS builds use different entrypoints. The enterprise entrypoint imports a
        // strict superset of the OSS entrypoint.
        app: [
            'react-hot-loader/patch',
            isEnterpriseBuild ? path.join(enterpriseDir, 'main.tsx') : path.join(__dirname, 'src', 'main.tsx'),
        ],

        'editor.worker': 'monaco-editor/esm/vs/editor/editor.worker.js',
        'json.worker': 'monaco-editor/esm/vs/language/json/json.worker',
    },
    output: {
        path: path.join(rootDir, 'ui', 'assets'),
        filename: 'scripts/[name].bundle.js',
        chunkFilename: 'scripts/[id]-[contenthash].chunk.js',
        publicPath: '/.assets/',
        globalObject: 'self',
        pathinfo: false,
    },
    devtool,
    plugins: [
        // Needed for React
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: JSON.stringify(mode),
            },
        }),
        new MiniCssExtractPlugin({ filename: 'styles/[name].bundle.css' }) as any, // @types package is incorrect
        new OptimizeCssAssetsPlugin(),
        new MonacoWebpackPlugin({
            languages: ['json'],
            features: [
                'bracketMatching',
                'clipboard',
                'coreCommands',
                'cursorUndo',
                'find',
                'format',
                'hover',
                'inPlaceReplace',
                'iPadShowKeyboard',
                'links',
                'suggest',
            ],
        }),
        new webpack.IgnorePlugin(/\.flow$/, /.*/),
    ],
    resolve: {
        extensions: ['.mjs', '.ts', '.tsx', '.js'],
        mainFields: ['es2015', 'module', 'browser', 'main'],
    },
    module: {
        rules: [
            // Run hot-loading-related Babel plugins on our application code only (because they'd be
            // slow to run on all JavaScript code).
            {
                test: /\.[jt]sx?$/,
                include: path.join(__dirname, 'src'),
                use: [
                    {
                        loader: 'babel-loader',
                        options: {
                            cacheDirectory: true,
                            plugins: [
                                'react-hot-loader/babel',
                                [
                                    '@sourcegraph/babel-plugin-transform-react-hot-loader-wrapper',
                                    {
                                        modulePattern: 'web/src/.*\\.tsx$',
                                        componentNamePattern: '(Page|Area)$',
                                    },
                                ],
                            ],
                        },
                    },
                ],
            },
            {
                test: /\.[jt]sx?$/,
                exclude: path.join(__dirname, 'src'),
                use: [
                    {
                        loader: 'babel-loader',
                        options: {
                            cacheDirectory: true,
                            configFile: path.join(__dirname, 'babel.config.js'),
                        },
                    },
                ],
            },
            {
                test: /\.(sass|scss)$/,
                use: [
                    mode === 'production' ? MiniCssExtractPlugin.loader : 'style-loader',
                    'css-loader',
                    {
                        loader: 'postcss-loader',
                        options: {
                            config: {
                                path: __dirname,
                            },
                        },
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            includePaths: [nodeModulesPath],
                        },
                    },
                ],
            },
            {
                // CSS rule for monaco-editor and other external plain CSS (skip SASS and PostCSS for build perf)
                test: /\.css$/,
                include: monacoEditorPaths,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },
}

export default config
