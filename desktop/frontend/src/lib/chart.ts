import { type ChartConfiguration, type Plugin } from "chart.js";

type Vertex = {
  x: number;
  y: number;
};

export function setupPlot(values: number[]): ChartConfiguration {
  const dataPoints: Vertex[] = [];

  const xValues = [];
  const yValues = [];

  for (let i = 0; i < values.length; i += 2) {
    const x = values[i];
    const y = values[i + 1];
    xValues.push(x);
    yValues.push(y);
    dataPoints.push({
      x,
      y,
    });
  }

  xValues.sort();
  yValues.sort();
  const xRange = xValues[xValues.length - 1] - xValues[0];
  const yRange = yValues[yValues.length - 1] - yValues[0];
  const marginRatio = 0.1;
  const xAxisMax = xValues[xValues.length - 1] + xRange * marginRatio;
  const xAxisMin = xValues[0] - xRange * marginRatio;
  const yAxisMax = yValues[yValues.length - 1] + yRange * marginRatio;
  const yAxisMin = yValues[0] - yRange * marginRatio;
  const data = {
    datasets: [
      {
        label: "Scatter Dataset",
        data: dataPoints,
        backgroundColor: "white",
        borderColor: "black",
      },
    ],
  };

  const plugin: Plugin = {
    id: "customCanvasBackgroundColor",
    beforeDraw: (chart, args, options) => {
      const { ctx } = chart;
      ctx.save();
      ctx.globalCompositeOperation = "destination-over";
      ctx.fillStyle = options.color || "#FFFFFF";
      ctx.fillRect(0, 0, chart.width, chart.height);
      ctx.restore();
    },
  };
  const config: ChartConfiguration = {
    type: "scatter",
    data: data,
    plugins: [plugin],
    options: {
      scales: {
        x: {
          max: xAxisMax,
          min: xAxisMin,
        },
        y: {
          max: yAxisMax,
          min: yAxisMin,
        },
      },
    },
  };
  return config;
}
