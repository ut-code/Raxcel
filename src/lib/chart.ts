import { type ChartConfiguration } from "chart.js";

type Vertex = {
  x: number;
  y: number;
}
export function setupPlot(values: number[]): ChartConfiguration {
  const rawData: Vertex[] = [];
  for (let i = 0; i < values.length / 2; i += 1) {
    rawData.push({
      x: values[i],
      y: values[i + values.length / 2]
    })
  }
  console.log(rawData);
  const data = {
    datasets: [{
      label: "Scatter Dataset",
      data: rawData,
      backgroundColor: "white",
      borderColor: "black",
    }]
  }

  const config: ChartConfiguration = {
    type: "scatter",
    data: data,
    options: {
      scales: {
        x: {
          type: "linear",
          position: "bottom"
        }
      }
    }
  }
  return config
}
