% Open the provided .fig file
openfig('ExecutionTime.fig');

% Load the data into MATLAB workspace
data = readtable('benchmark_v2_results.csv'); % Replace with your actual CSV if needed

% Extract unique iterations and prepare color gradient
uniqueIterations = unique(data.Iterations);
numColors = numel(uniqueIterations);
colors = jet(numColors); % Jet colormap for gradient

% Create the 2D plot
figure;
hold on;
for i = 1:numColors
    iter = uniqueIterations(i);
    iterData = data(data.Iterations == iter, :);
    plot(iterData.Goroutines, iterData.ExecutionTime_ms_, 'Color', colors(i, :), 'LineWidth', 1.5);
end
hold off;

% Add labels, legend, and title
xlabel('Number of Goroutines');
ylabel('Execution Time (ms)');
title('Execution Time vs Number of Goroutines');
colormap(colors);
colorbar('Ticks', linspace(0, 1, numColors), 'TickLabels', uniqueIterations);
legend(arrayfun(@num2str, uniqueIterations, 'UniformOutput', false), 'Location', 'bestoutside');
grid on;

% Save the figure as a .fig file
savefig('ExecutionTime_2D.fig');

disp('2D Execution Time figure saved as "ExecutionTime_2D.fig".');
