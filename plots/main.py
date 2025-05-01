import matplotlib.pyplot as plt
import pandas as pd
import os
import glob
import numpy as np
from math import pi

DATA_DIR = "./data"
COLORS = {
    "Go": "#01aed8",
    "Rust": "#e43715",
    "JS": "#f7df1d",
}

# Load CSV files
all_files = glob.glob(os.path.join(DATA_DIR, "*.csv"))
dfs = []
for f in all_files:
    df = pd.read_csv(f)
    df = df.drop(columns=[df.columns[0]])  # Drop first unnamed column
    dfs.append(df)

if not dfs:
    raise ValueError("No CSV files found in ./data")

df_all = pd.concat(dfs, ignore_index=True)
langs = df_all["lang"].unique()  # Go, Rust, JS
browsers = df_all["browser"].unique()  # Firefox, Chromium, GnomeWeb

# 1. FPS vs bunnies (line plot)
fig, axs = plt.subplots(
    len(langs), len(browsers), figsize=(15, 12), sharex=True, sharey=True
)
axs = np.atleast_2d(axs)
fig.suptitle("FPS vs Bunnies", fontsize=16)

for i, lang in enumerate(langs):
    for j, browser in enumerate(browsers):
        ax = axs[i, j]
        sub_df = df_all[(df_all["lang"] == lang) & (df_all["browser"] == browser)]
        if not sub_df.empty:
            ax.plot(
                sub_df["bunnies"], sub_df["fps_js"], label="Browser FPS", color="green"
            )
            ax.plot(
                sub_df["bunnies"],
                sub_df["fps_game"],
                label="Game FPS",
                color=COLORS.get(lang, "gray"),
            )
        ax.set_title(f"{lang} - {browser}")
        if i == len(langs) - 1:
            ax.set_xlabel("Bunnies")
        if j == 0:
            ax.set_ylabel("FPS")
        ax.legend()

plt.tight_layout(rect=[0, 0, 1, 0.96])
plt.show()

# 2. Max bunnies (bar chart)
max_bunnies_df = df_all.groupby(["lang", "browser"])["bunnies"].max().unstack()
max_bunnies_df.plot(kind="bar", figsize=(10, 6), colormap="Set2")
plt.title("Max Bunnies per Lang/Browser")
plt.ylabel("Max Bunnies")
plt.xlabel("Language")
plt.legend(title="Browser")
plt.tight_layout()
plt.show()

# 3. TPS stability
fig, axs = plt.subplots(
    len(langs), len(browsers), figsize=(15, 12), sharex=True, sharey=True
)
axs = np.atleast_2d(axs)
fig.suptitle("TPS Stability", fontsize=16)

for i, lang in enumerate(langs):
    for j, browser in enumerate(browsers):
        ax = axs[i, j]
        sub_df = df_all[(df_all["lang"] == lang) & (df_all["browser"] == browser)]
        if not sub_df.empty:
            ax.plot(sub_df["bunnies"], sub_df["tps"], marker="o", color="green")
        ax.set_title(f"{lang} - {browser}")
        if i == len(langs) - 1:
            ax.set_xlabel("Bunnies")
        if j == 0:
            ax.set_ylabel("TPS")

plt.tight_layout(rect=[0, 0, 1, 0.96])
plt.show()

# 4. Frame-time distribution (box plot)
frame_metrics = ["avg_frame", "min_frame", "max_frame"]
df_clean = df_all.dropna(subset=frame_metrics)

fig, axs = plt.subplots(1, 3, figsize=(18, 6))
# axs = np.atleast_2d(axs)
fig.suptitle("Frame Time Distribution by Language", fontsize=16)

for idx, metric in enumerate(frame_metrics):
    df_clean.boxplot(column=metric, by="lang", ax=axs[idx])
    axs[idx].set_title(f"{metric}")
    axs[idx].set_ylabel("ms")

plt.tight_layout(rect=[0, 0, 1, 0.95])
plt.show()

# 5. Heap usage over bunnies
fig, axs = plt.subplots(
    len(langs), len(browsers), figsize=(15, 12), sharex=True, sharey=True
)
axs = np.atleast_2d(axs)
fig.suptitle("Heap Memory Usage vs Bunnies", fontsize=16)

for i, lang in enumerate(langs):
    for j, browser in enumerate(browsers):
        ax = axs[i, j]
        sub_df = df_all[(df_all["lang"] == lang) & (df_all["browser"] == browser)]
        if not sub_df.empty:
            ax.plot(sub_df["bunnies"], sub_df["heap_mb"], marker="o", color="purple")
        ax.set_title(f"{lang} - {browser}")
        if i == len(langs) - 1:
            ax.set_xlabel("Bunnies")
        if j == 0:
            ax.set_ylabel("Heap (MB)")

plt.tight_layout(rect=[0, 0, 1, 0.96])
plt.show()

# 6. Click Latency
if "click_latency_ms" in df_all.columns:
    df_click = df_all.dropna(subset=["click_latency_ms"])
    # plt.figure(figsize=(8, 6))
    df_click.boxplot(column="click_latency_ms", by="lang")
    plt.title("Click Latency Distribution by Language")
    plt.ylabel("Latency (ms)")
    plt.tight_layout()
    plt.suptitle("")
    plt.show()

# 7. Radar chart: summary
summary_df = (
    df_all.groupby("lang")
    .agg(
        {
            "fps_js": "mean",
            "fps_game": "mean",
            "tps": "mean",
            "heap_mb": "mean",
            "click_latency_ms": "mean",
        }
    )
    .dropna()
)

summary_norm = (summary_df - summary_df.min()) / (summary_df.max() - summary_df.min())
categories = summary_norm.columns.tolist()
N = len(categories)
angles = [n / float(N) * 2 * pi for n in range(N)]
angles += angles[:1]

plt.figure(figsize=(8, 8))
for lang in summary_norm.index:
    values = summary_norm.loc[lang].tolist()
    values += values[:1]
    plt.polar(angles, values, label=lang, color=COLORS.get(lang, None))

plt.xticks(angles[:-1], categories)
plt.title("Performance Profile (Radar Chart)")
plt.legend(loc="upper right", bbox_to_anchor=(1.3, 1))
plt.tight_layout()
plt.show()
