% File paths
file1 = 'new_benchmark_results.csv';
file2 = 'new_new_benchmark_results.csv';
file3 = 'new_new_new_benchmark_results.csv';

% Read data from CSV files
data1 = readtable(file1, 'Delimiter', ';', 'VariableNamingRule', 'preserve');
data2 = readtable(file2, 'Delimiter', ';', 'VariableNamingRule', 'preserve');
data3 = readtable(file3, 'Delimiter', ';', 'VariableNamingRule', 'preserve');

% Concatenate data from all files
data = [data1; data2; data3];

% Extract variables
iterations = data.Iterations;
goroutines = data.Goroutines;
executionTime = data.("ExecutionTime(ms)");

% Create a grid for interpolation
[xq, yq] = meshgrid(linspace(min(iterations), max(iterations), 100), ...
                     linspace(min(goroutines), max(goroutines), 100));

% Interpolate execution time over the grid
zq = griddata(iterations, goroutines, executionTime, xq, yq, 'natural');

% Create a surface plot
figure;
surf(xq, yq, zq, 'EdgeColor', 'none'); % Smooth surface
xlabel('Iterations');
ylabel('Goroutines');
zlabel('Execution Time (ms)');
title('Surface Plot: Execution Time vs Iterations and Goroutines');
colorbar;
grid on;

% Save the figure as a .fig file
savefig('execution_time_3D_surface.fig');
