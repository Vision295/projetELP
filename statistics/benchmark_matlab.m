% Load data from the CSV file
data = readtable('benchmark_results.csv');

% Extract columns
iterations = data.Iterations;  % Number of iterations
goroutines = data.Goroutines;  % Number of goroutines
executionTime = data.ExecutionTime_ms_;  % Execution time in ms

% Reshape data for 3D plotting
% Assume the data is structured as a grid with regular steps
unique_iterations = unique(iterations);
unique_goroutines = unique(goroutines);

[X, Y] = meshgrid(unique_goroutines, unique_iterations);

% Reshape executionTime into a grid for plotting
Z = reshape(executionTime, length(unique_goroutines), length(unique_iterations))';

% Plot the data
figure;
surf(X, Y, Z);

% Add labels and title
xlabel('Number of Goroutines');
ylabel('Number of Iterations');
zlabel('Execution Time (ms)');
title('Execution Time vs Goroutines and Iterations');

% Add grid and color
grid on;
colormap('jet');
colorbar;
