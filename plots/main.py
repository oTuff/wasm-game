import glob
import os

import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


def main():
    # plot_fps_vs_bunnies()
    # plot_max_bunnies()
    # plot_tps_stability()
    # plot_frame_time_distribution()
    # plot_heap_usage()
    plot_load_time()


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
    opt_z_df_list = []
    opt_3_df_list = []

    for f in all_files:
        df = pd.read_csv(f)
        df = df.drop(columns=[df.columns[0]])

        if "laptop" in os.path.basename(f):
            laptop_df_list.append(df)
        elif "opt_z" in os.path.basename(f):
            opt_z_df_list.append(df)
        elif "opt_3" in os.path.basename(f):
            opt_3_df_list.append(df)
        else:
            df_list.append(df)

    if not df_list and not laptop_df_list:
        raise ValueError("No CSV files found in ./data")

    return df_list, laptop_df_list, opt_z_df_list, opt_3_df_list


df_list, laptop_df_list, opt_z, opt_3 = load_data()

df_all = pd.concat(df_list + opt_z, ignore_index=True)
# df_all = pd.concat(df_list + opt_3, ignore_index=True)
# df_all = pd.concat(laptop_df_list, ignore_index=True)

langs = df_all["lang"].unique()  # Go, Rust, JS
browsers = df_all["browser"].unique()  # Firefox, Chromium, GnomeWeb


# Function to plot FPS vs bunnies
def plot_fps_vs_bunnies():
    g = sns.FacetGrid(
        df_all,
        col="browser",
        hue="lang",
        palette=COLORS,
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
    df_clean = df_all.dropna(subset=["tps", "lang", "browser"])
    df_clean["lang_browser"] = df_clean["lang"] + " - " + df_clean["browser"]

    plt.figure(figsize=(12, 6))
    sns.boxplot(x="lang_browser", y="tps", data=df_clean)

    # plt.ylim(1, 70)  # Focused TPS range
    plt.xticks(rotation=45)
    plt.xlabel("Language - Browser")
    plt.ylabel("TPS")
    plt.tight_layout()

    finalize_plot("TPS_Stability")


# Function to plot Frame-time distribution (box plot)
def plot_frame_time_distribution():
    frame_metrics = ["avg_frame", "min_frame", "max_frame"]
    df_clean = df_all.dropna(subset=frame_metrics)

    fig, axs = plt.subplots(1, 3, figsize=(18, 6), sharey=True)

    for idx, metric in enumerate(frame_metrics):
        sns.boxplot(x="lang", y=metric, data=df_clean, ax=axs[idx])
        axs[idx].set_title(metric)
        axs[idx].set_ylabel("ms")

    plt.tight_layout(rect=[0, 0.03, 1, 0.95])  # Adjust layout to make room for suptitle

    finalize_plot("Frame_Time_Distribution")


# Function to plot Heap usage over bunnies
def plot_heap_usage():
    df_clean = df_all[df_all["browser"] == "Chromium"].dropna(subset=["heap_mb"])

    plt.figure(figsize=(8, 6))
    sns.lineplot(
        data=df_clean, x="bunnies", y="heap_mb", hue="lang", marker="o", palette=COLORS
    )

    plt.xlabel("Bunnies")
    plt.ylabel("Heap (MB)")
    plt.legend(title="Language")
    plt.tight_layout()

    finalize_plot("Heap_Usage")


def plot_load_time():
    max_bunnies = df_all.groupby(["lang", "browser"], as_index=False)[
        "wasm_exec_ms"
    ].max()
    sns.barplot(
        max_bunnies, x="lang", y="wasm_exec_ms", hue="browser", palette="icefire"
    )
    plt.ylabel("Load time (ms)")
    plt.xlabel("Language")

    finalize_plot("Load_Time")


# Helper function to either show or save the plot
def finalize_plot(plot_name):
    if SAVE_PLOTS_PNG:
        plt.savefig(f"{IMG_DIR}/{plot_name}.png")
        plt.close()
    else:
        plt.show()


if __name__ == "__main__":
    main()
