import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


SIZES = ["xsmall", "small", "medium", "large", "xlarge"]

def get_data():
    data = {
        "xsmall": {
            "sequential": 0.0,
            "2-thread": 0.0,
            "4-thread": 0.0,
            "6-thread": 0.0,
            "8-thread": 0.0,
            "12-thread": 0.0,
        }
    }

    for size in SIZES:
        data[size] = {}
        with open(f"./{size}/sequential.txt", "r") as f:
            times = f.readlines()
        data[size]["sequential"] = np.mean([float(time.strip()) for time in times])
        for i in [2, 4, 6, 8, 12]:
            with open(f"./{size}/{i}-thread.txt", "r") as f:
                times = f.readlines()
            data[size][f"{i}-thread"] = np.mean([float(time.strip()) for time in times])
    return data

def get_speedup(data):
    speedup = {}
    for size in SIZES:
        speedup[size] = {}
        for i in [2, 4, 6, 8, 12]:
            speedup[size][f"{i}-thread"] = data[size]["sequential"] / data[size][f"{i}-thread"]
    return speedup

def graph_speedup(speedup):
    for size in SIZES:
        x = [2, 4, 6, 8, 12]
        y = [speedup[size][f"{i}-thread"] for i in x]
        plt.plot(x, y, label=size)
    plt.xlabel("Number of threads")
    plt.ylabel("Speedup")
    plt.title("Speedup Graph")
    plt.legend()
    plt.savefig("my-speedup.png")
    print("Graph is saved successfully.")

if __name__ == "__main__":
    data = get_data()
    speedup = get_speedup(data)
    graph_speedup(speedup)