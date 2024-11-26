import os
import sys
import warnings
import numpy as np
import pandas as pd
import torch
from lightautoml.automl.presets.tabular_presets import TabularAutoML
import io
import pickle
import json

# Suppress all warnings and redirect stderr
sys.stderr = open(os.devnull, 'w')

N_THREADS = 4
N_FOLDS = 5
RANDOM_STATE = 42
TEST_SIZE = 0.2
TIMEOUT = 300
TARGET_NAME = '1'

def main():
    np.random.seed(RANDOM_STATE)
    torch.set_num_threads(N_THREADS)
    
    # Load model
    with open("internal/ml/models/model.pkl", "rb") as file:
        model = pickle.load(file)
    
    # Load data row
    with open("internal/ml/scripts/row.txt", "r") as file:
        data = file.readline().strip()
    
    # Prepare DataFrame
    df = pd.read_csv(io.StringIO(data), header=None)
    missing_values = [1, 2041, 7.31]
    missing_df = pd.DataFrame([missing_values])
    df = pd.concat([missing_df, df], axis=1).reset_index(drop=True)
    df.columns = pd.Index([str(i) for i in range(1, 91)])
    
    # Predict
    result = model.predict(df)
    
    # Extract prediction
    prediction = bool(result.data[0][0] > 0.5)
    probability = float(result.data[0][0])
    
    # Prepare output and save to JSON file
    output = {
        "prediction": prediction,
        "probability": probability
    }

    # Save result to a JSON file
    with open("internal/ml/scripts/prediction_result.json", "w") as json_file:
        json.dump(output, json_file)

if __name__ == "__main__":
    main()
