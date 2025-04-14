import matplotlib.pyplot as plt
import numpy as np

# Sample data
x = np.linspace(0, 10, 100)
y = np.sin(x)

# Create the plot
plt.plot(x, y, label="Sine Wave")
plt.xlabel("X-axis")
plt.ylabel("Y-axis")
plt.title("Sine Wave")
plt.legend()

# Show the plot
plt.show()
