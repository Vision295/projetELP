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
quality = data.Quality;

% Create a grid for interpolation
[xq, yq] = meshgrid(linspace(min(iterations), max(iterations), 100), ...
                     linspace(min(goroutines), max(goroutines), 100));

% Interpolate quality over the grid
zq = griddata(iterations, goroutines, quality, xq, yq, 'natural');

% Create a surface plot
figure;
surf(xq, yq, zq, 'EdgeColor', 'none'); % Smooth surface
xlabel('Iterations');
ylabel('Goroutines');
zlabel('Quality');
title('Surface Plot: Quality vs Iterations and Goroutines');
colorbar;
grid on;

% Save the figure as a .fig file
savefig('quality_3D_surface.fig');
