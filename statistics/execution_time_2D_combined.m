% Load the two data files
data1 = readtable("C:\Users\User\3TC\Projets\Période 2\ELP\projetELP\statistics\v2\benchmark_v2_results.csv"); % Replace with your first data file path
data2 = readtable("C:\Users\User\3TC\Projets\Période 2\ELP\projetELP\statistics\benchmark_v3_results.csv"); % Replace with your second data file path

% Concatenate the two datasets
combinedData = [data1; data2];

% Extract unique iterations and prepare color gradient
uniqueIterations = unique(combinedData.Iterations);
numColors = numel(uniqueIterations);
colors = jet(numColors); % Jet colormap for gradient

% Create the 2D plot
figure;
hold on;
for i = 1:numColors
    iter = uniqueIterations(i);
    iterData = combinedData(combinedData.Iterations == iter, :);
    plot(iterData.Goroutines, iterData.ExecutionTime_ms_, 'Color', colors(i, :), 'LineWidth', 1.5);
end
hold off;

% Add labels, legend, and title
xlabel('Number of Goroutines');
ylabel('Execution Time (ms)');
title('Execution Time (Combined Data)');
colormap(colors);
colorbar('Ticks', linspace(0, 1, numColors), 'TickLabels', uniqueIterations);
legend(arrayfun(@num2str, uniqueIterations, 'UniformOutput', false), 'Location', 'bestoutside');
grid on;

% Save the combined figure as a .fig file
savefig('ExecutionTime_Combined_2D.fig');

disp('Combined 2D Execution Time figure saved as "ExecutionTime_Combined_2D.fig".');
