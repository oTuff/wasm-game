use bevy::{
    diagnostic::{DiagnosticsStore, FrameTimeDiagnosticsPlugin},
    prelude::*,
};
use fastrand;
use wasm_bindgen::prelude::*;

const SCREEN_WIDTH: f32 = 640.0;
const SCREEN_HEIGHT: f32 = 480.0;
const UPPER_BOUND: i32 = 40; // Amount of pixels from the top
const BUNNY_SCALE: f32 = 0.2;

#[wasm_bindgen]
pub fn run_bevy_app() {
    App::new()
        // Set window resolution
        .add_plugins((
            DefaultPlugins.set(WindowPlugin {
                primary_window: Some(Window {
                    resolution: (SCREEN_WIDTH, SCREEN_HEIGHT).into(),
                    ..default()
                }),
                ..default()
            }),
            FrameTimeDiagnosticsPlugin::default(),
        ))
        .add_systems(Startup, setup)
        .add_systems(
            Update,
            (update_fps_text, initialize_bunny_image_data, spawn_bunny),
        )
        .run();
}

// marker component
#[derive(Component)]
struct ColorText;

// marker component
#[derive(Component, Default)]
struct Bunny;

// Instead of creating separate Position and Scale components,
// use the `Transform` inside `SpriteBundle` to store position and scale.
//
// #[derive(Component, Default)]
// struct Position {
//     x: f32,
//     y: f32,
// }
// #[derive(Component, Default)]
// struct Scale {
//     x: f32,
//     y: f32,
// }

#[derive(Component, Default)]
struct Velocity {
    x: f32,
    y: f32,
}

#[derive(Bundle, Default)]
struct BunnyBundle {
    bunny: Bunny,
    sprite_bundle: SpriteBundle,
    velocity: Velocity,
}

#[derive(Resource)]
struct BunnyImage {
    handle: Handle<Image>,
    width: f32,
    height: f32,
}

fn initialize_bunny_image_data(
    mut my_image: ResMut<BunnyImage>,
    images: Res<Assets<Image>>,
    mut ev_asset: EventReader<AssetEvent<Image>>,
) {
    if my_image.width != 0.0 {
        info!("not running anymore");
        return;
    }
    for event in ev_asset.read() {
        match event {
            AssetEvent::Added { id } => {
                info!("Asset added with id: {:?}", id);

                if let Some(image) = images.get(*id) {
                    let size = image.size_f32();
                    my_image.width = size.x;
                    my_image.height = size.y;
                    info!("Loaded image size: {} x {}", size.x, size.y);
                }
            }
            AssetEvent::LoadedWithDependencies { id } => {
                info!("Asset loaded with dependencies with id: {:?}", id);
                if let Some(image) = images.get(*id) {
                    let size = image.size_f32();
                    my_image.width = size.x;
                    my_image.height = size.y;
                    info!(
                        "Loaded(with dependencies) image size: {} x {}",
                        size.x, size.y
                    );
                }
            }
            AssetEvent::Modified { id } => {
                info!("Asset modified with id: {:?}", id);
            }
            AssetEvent::Removed { id } => {
                info!("Asset removed with id: {:?}", id);
            }
            AssetEvent::Unused { id } => {
                info!(
                    "Asset unused (last strong handle dropped) with id: {:?}",
                    id
                );
            }
        }
    }
}

fn spawn_bunny(mut commands: Commands, my_image: Res<BunnyImage>) {
    let texture_handle = my_image.handle.clone();

    let position = Vec3::new(
        -SCREEN_WIDTH / 2.0 + (my_image.width / 2.0 * BUNNY_SCALE),
        SCREEN_HEIGHT / 2.0 - (my_image.height / 2.0 * BUNNY_SCALE),
        0.0,
    );

    let velocity = Velocity {
        x: fastrand::f32() * 2.0 + 2.0,
        y: fastrand::f32() * 2.0 + 2.0,
    };

    commands.spawn(BunnyBundle {
        sprite_bundle: SpriteBundle {
            sprite: texture_handle.into(),
            transform: Transform {
                translation: position,
                scale: Vec3::splat(BUNNY_SCALE),
                ..default()
            },
            ..default()
        },
        velocity,
        ..default()
    });
}

fn setup(mut commands: Commands, asset_server: Res<AssetServer>) {
    // UI camera
    commands.spawn(Camera2d);

    // load image
    let handle = asset_server.load("bunny.png");
    commands.insert_resource(BunnyImage {
        handle,
        width: 0.0,
        height: 0.0,
    });

    // Debug text
    commands.spawn((Text::new("FPS:"), ColorText));
}

// Debug information
fn update_fps_text(
    diagnostics: Res<DiagnosticsStore>,
    mut writer: TextUiWriter,
    text_query: Query<Entity, With<ColorText>>,
) {
    if let Some(fps_diag) = diagnostics.get(&FrameTimeDiagnosticsPlugin::FPS) {
        if let Some(fps) = fps_diag.smoothed() {
            for entity in &text_query {
                *writer.text(entity, 0) = format!("FPS: {:.1}", fps);
            }
        }
    }
}
