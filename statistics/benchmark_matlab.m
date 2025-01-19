% Load the CSV file
data = readtable('benchmark_v2_results.csv', 'Delimiter', ';');

% Extract columns from the table
iterations = data.Iterations;
goroutines = data.Goroutines;
executionTime = data.ExecutionTime_ms_;
ssimScore = data.SSIM_Score;

% Ensure data types are correct
iterations = double(iterations);
goroutines = double(goroutines);
executionTime = double(executionTime);
ssimScore = double(ssimScore);

% Unique values for iterations and goroutines
uniqueIterations = unique(iterations);
uniqueGoroutines = unique(goroutines);

% Create a grid for iterations and goroutines
[iterGrid, goroutinesGrid] = meshgrid(uniqueIterations, uniqueGoroutines);

% Reshape executionTime and ssimScore into matrices for surf plotting
executionTimeMatrix = reshape(executionTime, length(uniqueGoroutines), length(uniqueIterations));
ssimScoreMatrix = reshape(ssimScore, length(uniqueGoroutines), length(uniqueIterations));
timeQualityMatrix = executionTimeMatrix .* ssimScoreMatrix;

% Plot 1: Execution Time
figure;
surf(iterGrid, goroutinesGrid, executionTimeMatrix);
xlabel('Number of Iterations');
ylabel('Number of Goroutines');
zlabel('Execution Time (ms)');
title('Execution Time');
colorbar;
view(3);

% Plot 2: SSIM Score
figure;
surf(iterGrid, goroutinesGrid, ssimScoreMatrix);
xlabel('Number of Iterations');
ylabel('Number of Goroutines');
zlabel('Image Quality (SSIM Score)');
title('Image Quality');
colorbar;
view(3);

% Plot 3: Time * Quality
figure;
surf(iterGrid, goroutinesGrid, timeQualityMatrix);
xlabel('Number of Iterations');
ylabel('Number of Goroutines');
zlabel('Execution Time * Image Quality');
title('Execution Time * Image Quality');
colorbar;
view(3);

% Save figures to disk
saveas(1, 'execution_time_vs_iterations_goroutines.png');
saveas(2, 'image_quality_vs_iterations_goroutines.png');
saveas(3, 'time_quality_vs_iterations_goroutines.png');
disp('Plots saved as PNG files.');

