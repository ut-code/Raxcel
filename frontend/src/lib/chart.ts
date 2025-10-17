import { type ChartConfiguration } from "chart.js";

type Vertex = {
  x: number;
  y: number;
};
export function setupPlot(values: number[]): ChartConfiguration {
  const rawData: Vertex[] = [];
  console.log(values);
  for (let i = 0; i < values.length; i += 2) {
    rawData.push({
      x: values[i],
      y: values[i + 1],
    });
  }
  const data = {
    datasets: [
      {
        label: "Scatter Dataset",
        data: rawData,
        backgroundColor: "white",
        borderColor: "black",
      },
    ],
  };

  const config: ChartConfiguration = {
    type: "scatter",
    data: data,
    options: {
      scales: {
        x: {
          type: "linear",
          position: "bottom",
        },
      },
    },
  };
  return config;
}
