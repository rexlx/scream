function createLineChart(id, data, label = "Line Chart") {
  const ctx = document.getElementById(canvasId).getContext('2d');

  // Extract labels and datasets from the data object
  const labels = Array.from({ length: data['room_created'].length }, (_, i) => i + 1);
  const datasets = Object.keys(data).filter(key => key !== 'room_created').map(key => ({
      label: key,
      data: data[key],
      borderColor: getRandomColor(), // Assign a random color to each line
      tension: 0.1
  }));

  new Chart(ctx, {
      type: 'line',
      data: {
          labels: labels,
          datasets: datasets
      },
      options: {
          scales: {
              y: {
                  beginAtZero: true
              }
          }
      }
  });
  return chart;
}