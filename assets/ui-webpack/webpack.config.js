const path = require('path');
const {CleanWebpackPlugin} = require('clean-webpack-plugin');
const MiniCSSExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const TerserJSPlugin = require('terser-webpack-plugin');
const webpack = require('webpack');

module.exports = (env, options) => {
    const config = {
        entry: {
            admin: "./src/js/admin/App.js",
            gallery: "./src/js/gallery/App.js"
        },
        output: {
            filename: '[name].bundle.js',
            path: path.resolve(__dirname, 'dist')
        },
        module: {
            rules: [
                {
                    test: /\.js$/,
                    exclude: /node_modules/,
                    use: {
                        loader: "babel-loader",
                        options: {
                            presets: [
                                "@babel/preset-env",
                                "@babel/preset-react"
                            ]
                        }
                    }
                },
                {
                    test: /\.(sa|sc|c)ss$/,
                    use: [
                        MiniCSSExtractPlugin.loader,
                        'css-loader',
                        'sass-loader'
                    ]
                }, {
                    test: /\.png$/,
                    loader: 'file-loader'
                }, {
                    test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                    loader: "url-loader?limit=10000&mimetype=application/font-woff"
                }, {
                    test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                    loader: "file-loader"
                }, {
                    test: /\.html$/,
                    include: [
                        path.resolve(__dirname, "src/html")
                    ],
                    loader: "html-loader",
                    options: {
                        minimize: true
                    }
                }, {
                    test: /\.txt$/i,
                    loader: 'raw-loader',
                    options: {
                        esModule: false
                    }
                },

            ]
        },
        plugins: [
            new CleanWebpackPlugin({}),
            new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
            new MiniCSSExtractPlugin({
                filename: '[name].bundle.css'
            })
        ]
    };

    if (options.mode === "production") {
        config.optimization = {
            minimizer: [
                new TerserJSPlugin({}),
                new OptimizeCSSAssetsPlugin({})
            ]
        }
    }

    return config;
};
