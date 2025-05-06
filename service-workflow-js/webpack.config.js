import path from "path";
import { fileURLToPath } from "url";
import nodeExternals from "webpack-node-externals";
import webpack from "webpack";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default {
  target: "node",
  mode: process.env.NODE_ENV === "production" ? "production" : "development",
  entry: {
    index: ["./src/index.ts"],
    worker: ["./src/worker.ts"],
    "kafka-consumer": ["./src/kafka-consumer.ts"],
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "[name].js",
    clean: true,
  },
  module: {
    rules: [
      {
        test: /\.(ts|js)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: "babel-loader",
            options: {
              presets: [
                [
                  "@babel/preset-env",
                  {
                    targets: {
                      node: "18",
                    },
                  },
                ],
                "@babel/preset-typescript",
              ],
            },
          },
          {
            loader: "ts-loader",
            options: {
              transpileOnly: true,
            },
          },
        ],
      },
    ],
  },
  externals: [nodeExternals()],
  resolve: {
    extensions: [".ts", ".js", ".json"],
    mainFiles: ["index"],
    fullySpecified: false,
  },
  plugins: [new webpack.HotModuleReplacementPlugin()],
  watch: process.env.NODE_ENV === "development",
  watchOptions: {
    ignored: /node_modules/,
    aggregateTimeout: 300,
  },
};
