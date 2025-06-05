import glob
import os

import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


def main():
    plot_fps_vs_bunnies()
    # plot_max_bunnies()
    # plot_tps_stability()
    # plot_frame_time_distribution()
    # plot_heap_usage()
    # plot_load_time()
    # rust_vs_rust()


DATA_DIR = "./data"
IMG_DIR = "./img"
LANG_COLORS = {
    "Go": "#01aed8",
    "Rust": "#e43715",
    "JS": "#f7df1d",
}
LANG_CONTRAST_COLORS = {
    "Go": "#FF1493",
    "Rust": "#00C2A0",
    "JS": "#3F51B5",
}
BROWSER_COLORS = {
    "Chromium": "#264AFF",
    "Firefox": "#D32F2F",
    "GnomeWeb": "#7E57C2",
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

# df_all = pd.concat(df_list + opt_z, ignore_index=True)
df_all = pd.concat(df_list + opt_3, ignore_index=True)
# df_all = pd.concat(df_list + opt_3 + opt_z, ignore_index=True)
# df_all = pd.concat(laptop_df_list, ignore_index=True)

langs = df_all["lang"].unique()  # Go, Rust, JS
browsers = df_all["browser"].unique()  # Firefox, Chromium, GnomeWeb

# rust vs rust performance
for df in opt_z:
    df["opt_level"] = "opt_z"

for df in opt_3:
    df["opt_level"] = "opt_3"

opt_combined = pd.concat(opt_z + opt_3, ignore_index=True)

opt_filtered = opt_combined[
    (opt_combined["lang"] == "Rust") & (opt_combined["browser"] == "Chromium")
]


def rust_vs_rust():
    plt.figure(figsize=(8, 6))

    sns.lineplot(
        data=opt_filtered,
        x="bunnies",
        y="fps_js",
        hue="opt_level",
        linewidth=3,
    )

    plt.axhline(y=60, color="red", linestyle="--", linewidth=1, label="60 FPS")

    for opt_level in opt_filtered["opt_level"].unique():
        df_opt = opt_filtered[opt_filtered["opt_level"] == opt_level].sort_values(
            "bunnies"
        )
        bunnies = df_opt["bunnies"].values
        fps = df_opt["fps_game"].values

        for i in range(len(fps) - 1):
            x0, x1 = bunnies[i], bunnies[i + 1]
            y0, y1 = fps[i], fps[i + 1]

            if x0 >= 1500 and x1 >= 1500 and (y0 - 60) * (y1 - 60) < 0:
                bunny_cross = x0 + (60 - y0) * (x1 - x0) / (y1 - y0)
                plt.plot(
                    bunny_cross,
                    60,
                    "X",
                    markersize=10,
                    label=f"{opt_level} {x0}",
                )
                break

    plt.xlabel("Bunnies")
    plt.ylabel("FPS (Game Loop)")
    plt.legend(title="Optimization Level")
    plt.tight_layout()

    finalize_plot("rust_vs_rust")


# Function to plot FPS vs bunnies
def plot_fps_vs_bunnies():
    for browser in browsers:
        print(browser)
        df_subset = df_all[df_all["browser"] == browser]

        plt.figure(figsize=(8, 6))
        sns.lineplot(
            data=df_subset,
            x="bunnies",
            y="fps_js",
            hue="lang",
            palette=LANG_COLORS,
            linewidth=3,
        )

        plt.axhline(y=60, color="red", linestyle="--", linewidth=1, label="60 FPS")

        for lang in df_subset["lang"].unique():
            df_lang = df_subset[df_subset["lang"] == lang].sort_values("bunnies")
            bunnies = df_lang["bunnies"].values
            fps = df_lang["fps_game"].values

            for i in range(len(fps) - 1):
                x0, x1 = bunnies[i], bunnies[i + 1]
                y0, y1 = fps[i], fps[i + 1]

                if x0 >= 1500 and x1 >= 1500 and (y0 - 60) * (y1 - 60) < 0:
                    bunny_cross = x0 + (60 - y0) * (x1 - x0) / (y1 - y0)
                    plt.plot(
                        bunny_cross,
                        60,
                        "X",
                        markersize=10,
                        color=LANG_CONTRAST_COLORS[lang],
                        label=f"{lang} {x0}",
                    )
                    break

        plt.xlabel("Bunnies")
        plt.ylabel("FPS")
        plt.legend(title="Language")
        plt.tight_layout()

        finalize_plot(f"FPS_vs_Bunnies_{browser}")


# Function to plot Max bunnies (bar chart)
def plot_max_bunnies():
    max_bunnies = df_all.groupby(["lang", "browser"], as_index=False)["bunnies"].max()

    ax = sns.barplot(
        max_bunnies,
        x="lang",
        y="bunnies",
        hue="browser",
        palette=BROWSER_COLORS,
        order=["Go", "Rust", "JS"],
    )

    for container in ax.containers:
        ax.bar_label(container, fmt="%.0f", label_type="edge")

    plt.ylabel("Max Bunnies")
    plt.xlabel("Language")

    finalize_plot("Max_Bunnies")


# Function to plot TPS stability
def plot_tps_stability():
    df_clean = df_all.dropna(subset=["tps", "lang", "browser"])
    df_clean["lang_browser"] = df_clean["lang"] + " - " + df_clean["browser"]

    plt.figure(figsize=(12, 6))
    sns.boxplot(x="lang_browser", y="tps", data=df_clean)

    plt.ylim(55, 65)
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

    plt.tight_layout(rect=[0, 0.03, 1, 0.95])

    finalize_plot("Frame_Time_Distribution")


# Function to plot Heap usage over bunnies
def plot_heap_usage():
    df_clean = df_all[df_all["browser"] == "Chromium"].dropna(subset=["heap_mb"])

    plt.figure(figsize=(8, 6))
    sns.lineplot(
        data=df_clean,
        x="bunnies",
        y="heap_mb",
        hue="lang",
        hue_order=["JS", "Rust", "Go"],
        palette=LANG_COLORS,
        linewidth=3,
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
