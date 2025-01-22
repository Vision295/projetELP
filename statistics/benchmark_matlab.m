% Load the CSV file
data = readtable('benchmark_v2_results.csv');

% Extract data from the table
Iterations = data.Iterations;
Goroutines = data.Goroutines;
ExecutionTime = data.ExecutionTime_ms_;
Quality = data.Quality;

% Compute ExecutionTime * Quality
ExecTime_Quality = ExecutionTime .* Quality;

% Reshape the data for 3D plots
[X, Y] = meshgrid(unique(Iterations), unique(Goroutines));
Z1 = griddata(Iterations, Goroutines, ExecutionTime, X, Y);
Z2 = griddata(Iterations, Goroutines, Quality, X, Y);
Z3 = griddata(Iterations, Goroutines, ExecTime_Quality, X, Y);

% Plot Execution Time vs Iterations and Goroutines
figure;
surf(X, Y, Z1);
xlabel('Iterations');
ylabel('Goroutines');
zlabel('Execution Time (ms)');
title('Execution Time');
savefig('ExecutionTime');

% Plot Quality vs Iterations and Goroutines
figure;
surf(X, Y, Z2);
xlabel('Iterations');
ylabel('Goroutines');
zlabel('Quality');
title('Quality');
savefig('Quality');

% Plot Execution Time * Quality vs Iterations and Goroutines
figure;
surf(X, Y, Z3);
xlabel('Iterations');
ylabel('Goroutines');
zlabel('Execution Time * Quality');
title('Execution Time * Quality');
savefig('ExecTime_Quality');

disp('Figures saved as MATLAB .fig files.');



