import glob
import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns


def main():
    # plot_fps_vs_bunnies()
    # plot_max_bunnies()
    # plot_tps_stability()
    # plot_frame_time_distribution()
    plot_heap_usage()
    # plot_click_latency()  # not relevant for js
    # plot_radar_chart()


DATA_DIR = "./data"
IMG_DIR = "./img"
COLORS = {
    "Go": "#01aed8",
    "Rust": "#e43715",
    "JS": "#f7df1d",
}
SAVE_PLOTS_PNG = False


# Load CSV files
def load_data():
    all_files = glob.glob(os.path.join(DATA_DIR, "*.csv"))

    df_list = []
    laptop_df_list = []

    for f in all_files:
        df = pd.read_csv(f)
        df = df.drop(columns=[df.columns[0]])

        if "laptop" in os.path.basename(f):
            laptop_df_list.append(df)
        else:
            df_list.append(df)

    if not df_list and not laptop_df_list:
        raise ValueError("No CSV files found in ./data")

    return df_list, laptop_df_list


df_list, laptop_df_list = load_data()
df_all = pd.concat(df_list, ignore_index=True)
# df_laptop = pd.concat(laptop_df_list, ignore_index=True)
# df_all = pd.concat(laptop_df_list, ignore_index=True)
langs = df_all["lang"].unique()  # Go, Rust, JS
browsers = df_all["browser"].unique()  # Firefox, Chromium, GnomeWeb


# Function to plot FPS vs bunnies
def plot_fps_vs_bunnies():
    g = sns.FacetGrid(
        df_all,
        col="browser",
        hue="lang",
        # palette=COLORS,
        col_wrap=len(browsers),
        height=6,
        aspect=1.5,
    )
    g.map(sns.lineplot, "bunnies", "fps_game", alpha=0.7)

    g.set_axis_labels("Bunnies", "FPS")
    g.set_titles("{col_name}")
    g.add_legend(title="Language")
    g.tight_layout(pad=2.0)

    finalize_plot("FPS_vs_Bunnies")


# Function to plot Max bunnies (bar chart)
def plot_max_bunnies():
    max_bunnies = df_all.groupby(["lang", "browser"], as_index=False)["bunnies"].max()
    sns.barplot(max_bunnies, x="lang", y="bunnies", hue="browser", palette="icefire")
    plt.ylabel("Max Bunnies")
    plt.xlabel("Language")

    finalize_plot("Max_Bunnies")


# Function to plot TPS stability
def plot_tps_stability():
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

    finalize_plot("TPS_Stability")


# Function to plot Frame-time distribution (box plot)
def plot_frame_time_distribution():
    frame_metrics = ["avg_frame", "min_frame", "max_frame"]
    df_clean = df_all.dropna(subset=frame_metrics)

    fig, axs = plt.subplots(1, 3, figsize=(18, 6))
    fig.suptitle("Frame Time Distribution by Language", fontsize=16)

    for idx, metric in enumerate(frame_metrics):
        df_clean.boxplot(column=metric, by="lang", ax=axs[idx])
        axs[idx].set_title(f"{metric}")
        axs[idx].set_ylabel("ms")

    finalize_plot("Frame_Time_Distribution")


# Function to plot Heap usage over bunnies
def plot_heap_usage():
    fig, ax = plt.subplots(figsize=(8, 6))

    for lang in langs:
        sub_df = df_all[(df_all["lang"] == lang) & (df_all["browser"] == "Chromium")]
        if not sub_df.empty:
            ax.plot(
                sub_df["bunnies"],
                sub_df["heap_mb"],
                marker="o",
                label=lang,
                color=COLORS.get(lang, "gray"),
            )

    ax.set_xlabel("Bunnies")
    ax.set_ylabel("Heap (MB)")
    ax.legend(title="Language")

    finalize_plot("Heap_Usage")


# Function to plot Click Latency (box plot)
def plot_click_latency():
    if "click_latency_ms" in df_all.columns:
        df_click = df_all.dropna(subset=["click_latency_ms"])
        df_click.boxplot(column="click_latency_ms", by="lang")
        plt.title("Click Latency Distribution by Language")
        plt.ylabel("Latency (ms)")
        plt.tight_layout()
        plt.suptitle("")
        finalize_plot("Click_Latency")


# Function to plot Radar chart (summary)

# def plot_radar_chart():
#     summary_df = (
#         df_all.groupby("lang")
#         .agg(
#             {
#                 "fps_js": "mean",
#                 "fps_game": "mean",
#                 "tps": "mean",
#                 "heap_mb": "mean",
#                 "click_latency_ms": "mean",
#             }
#         )
#         .dropna()
#     )
#
#     summary_norm = (summary_df - summary_df.min()) / (
#         summary_df.max() - summary_df.min()
#     )
#     categories = summary_norm.columns.tolist()
#     N = len(categories)
#     angles = [n / float(N) * 2 * pi for n in range(N)]
#     angles += angles[:1]
#
#     plt.figure(figsize=(8, 8))
#     for lang in summary_norm.index:
#         values = summary_norm.loc[lang].tolist()
#         values += values[:1]
#         plt.polar(angles, values, label=lang, color=COLORS.get(lang, None))
#
#     plt.xticks(angles[:-1], categories)
#     plt.title("Performance Profile (Radar Chart)")
#     plt.legend(loc="upper right", bbox_to_anchor=(1.3, 1))
#     finalize_plot("Radar_Chart")


# Helper function to either show or save the plot
def finalize_plot(plot_name):
    if SAVE_PLOTS_PNG:
        plt.savefig(f"{IMG_DIR}/{plot_name}.png")
        plt.close()
    else:
        plt.show()


if __name__ == "__main__":
    main()
