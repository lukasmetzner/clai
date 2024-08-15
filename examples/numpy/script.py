import numpy as np

# Set a seed for reproducibility
np.random.seed(42)

# Generate a random 4x4 matrix
matrix = np.random.rand(4, 4)

# Generate a random 4x1 vector
vector = np.random.rand(4, 1)

# Perform a matrix-vector multiplication
result = np.dot(matrix, vector)

# Calculate the mean of the result
mean_result = np.mean(result)

print("Random 4x4 Matrix:\n", matrix)
print("\nRandom 4x1 Vector:\n", vector)
print("\nMatrix-Vector Multiplication Result:\n", result)
print("\nMean of the Result:", mean_result)

