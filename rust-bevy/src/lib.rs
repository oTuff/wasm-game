use bevy::{
    asset::AssetMetaCheck,
    diagnostic::{DiagnosticsStore, FrameTimeDiagnosticsPlugin},
    prelude::*,
    winit::{UpdateMode, WinitSettings},
};
use fastrand;
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use std::sync::RwLock;
use wasm_bindgen::prelude::*;

const SCREEN_WIDTH: f32 = 640.0;
const SCREEN_HEIGHT: f32 = 480.0;
const UPPER_BOUND: f32 = 40.0; // Amount of pixels from the top
const BUNNY_SCALE: f32 = 0.2;

#[wasm_bindgen]
pub fn run_bevy_app() {
    App::new()
        .add_plugins((
            DefaultPlugins
                // set window resolution
                .set(WindowPlugin {
                    primary_window: Some(Window {
                        resolution: (SCREEN_WIDTH, SCREEN_HEIGHT).into(),
                        ..default()
                    }),
                    ..default()
                })
                // don't look for `.meta` files
                .set(AssetPlugin {
                    meta_check: AssetMetaCheck::Never,
                    ..default()
                }),
            FrameTimeDiagnosticsPlugin::default(),
        ))
        // fix framerate in wasm when unfocused
        .insert_resource(WinitSettings {
            focused_mode: UpdateMode::Continuous,
            unfocused_mode: UpdateMode::Continuous,
        })
        // set background color
        .insert_resource(ClearColor(Color::BLACK))
        .add_systems(Startup, setup)
        .add_systems(Update, (initialize_bunny_image_data, input_system))
        .add_systems(
            FixedUpdate,
            (update_positions, check_boundaries, update_diagnostics_text),
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
    sprite: Sprite,
    transform: Transform,
    velocity: Velocity,
}

#[derive(Resource)]
struct BunnyImage {
    handle: Handle<Image>,
    width: f32,
    height: f32,
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
    commands.spawn((
        Text::new(""),
        TextFont {
            font_size: 11.0,
            ..default()
        },
        ColorText,
    ));
}

// Get the image size
fn initialize_bunny_image_data(
    commands: Commands,
    mut bunny_image: ResMut<BunnyImage>,
    images: Res<Assets<Image>>,
) {
    // stop running once the width(and height) is set by this system
    if bunny_image.width != 0.0 {
        return;
    }
    // get the width and height of the image
    if let Some(image) = images.get(&mut bunny_image.handle) {
        let size = image.size_f32();
        bunny_image.width = size.x;
        bunny_image.height = size.y;
        spawn_bunnies(10, commands, bunny_image.into());
    }
}

fn spawn_bunnies(amount: u32, mut commands: Commands, bunny_image: Res<BunnyImage>) {
    for _ in 0..amount {
        // set the initial position
        let position = Vec3::new(
            -SCREEN_WIDTH / 2.0 + (bunny_image.width / 2.0 * BUNNY_SCALE) + fastrand::f32() * 5.0,
            SCREEN_HEIGHT / 2.0 - (bunny_image.height / 2.0 * BUNNY_SCALE) - fastrand::f32() * 5.0,
            0.0,
        );

        let velocity = Velocity {
            x: fastrand::f32() * 2.0 + 2.0,
            y: fastrand::f32() * 2.0 + 2.0,
        };

        commands.spawn(BunnyBundle {
            sprite: Sprite {
                image: bunny_image.handle.clone(),
                custom_size: Some(Vec2::new(
                    bunny_image.width * BUNNY_SCALE,
                    bunny_image.height * BUNNY_SCALE,
                )),
                ..default()
            },
            transform: Transform {
                translation: position,
                ..default()
            },
            velocity,
            ..default()
        });
    }
}

fn update_positions(mut bunny_query: Query<(&mut Velocity, &mut Transform), With<Bunny>>) {
    for (mut velocity, mut transform) in bunny_query.iter_mut() {
        transform.translation.x += velocity.x;
        transform.translation.y += velocity.y;
        velocity.y -= 0.75; // hardcoded gravity
    }
}

fn check_boundaries(
    mut bunny_query: Query<(&mut Velocity, &mut Transform), With<Bunny>>,
    bunny_image: ResMut<BunnyImage>,
) {
    let scaled_width = bunny_image.width / 2.0 * BUNNY_SCALE;
    let scaled_height = bunny_image.height / 2.0 * BUNNY_SCALE;

    for (mut velocity, mut transform) in bunny_query.iter_mut() {
        if transform.translation.x > SCREEN_WIDTH / 2.0 - scaled_width
            || transform.translation.x < -SCREEN_WIDTH / 2.0 + scaled_width
        {
            velocity.x = -velocity.x;
        }

        if transform.translation.y < -SCREEN_HEIGHT / 2.0 + scaled_height {
            transform.translation.y = -SCREEN_HEIGHT / 2.0 + scaled_height;
            velocity.y = -velocity.y;
        }
        if transform.translation.y > SCREEN_HEIGHT / 2.0 - UPPER_BOUND - scaled_height
            && velocity.y > 0.0
        {
            velocity.y *= 0.7;
        }
    }
}

fn input_system(
    mouse_button_input: Res<ButtonInput<MouseButton>>,
    commands: Commands,
    bunny_image: ResMut<BunnyImage>,
    bunnies: Query<(), With<Bunny>>,
) {
    let bunny_count = bunnies.iter().count() as u32;

    if mouse_button_input.pressed(MouseButton::Left) {
        spawn_bunnies(1, commands, bunny_image.into());
    } else if mouse_button_input.just_pressed(MouseButton::Middle) {
        let mut to_add = 1000 - (bunny_count % 1000);
        if to_add == 0 {
            to_add = 1000;
        }
        spawn_bunnies(to_add, commands, bunny_image.into());
    } else if mouse_button_input.just_pressed(MouseButton::Right) {
        let mut to_add = 100 - (bunny_count % 100);
        if to_add == 0 {
            to_add = 100;
        }
        spawn_bunnies(to_add, commands, bunny_image.into());
    }
}

#[derive(Serialize, Deserialize, Default)]
pub struct Metrics {
    pub fps: f64,
    pub tps: f64,
    pub bunnies: usize,
}

// Lazily initializes `RwLock<Metrics>` for thread-safe access to metrics
static GLOBAL_METRICS: Lazy<RwLock<Metrics>> = Lazy::new(|| RwLock::new(Metrics::default()));

// Debug information
fn update_diagnostics_text(
    diagnostics: Res<DiagnosticsStore>,
    mut writer: TextUiWriter,
    text_query: Query<Entity, With<ColorText>>,
    fixed_time: Res<Time<Fixed>>,
    bunnies: Query<(&mut Velocity, &mut Transform), With<Bunny>>,
) {
    let bunny_count = bunnies.iter().count();
    let fixed_delta = fixed_time.delta_secs_f64();
    let tickrate = if fixed_delta > 0.0 {
        1.0 / fixed_delta
    } else {
        0.0
    };
    if let Some(fps_diag) = diagnostics.get(&FrameTimeDiagnosticsPlugin::FPS) {
        if let Some(fps) = fps_diag.smoothed() {
            for entity in &text_query {
                *writer.text(entity, 0) = format!(
                    "FPS: {:.0}\nTPS: {:.0}\nbunnies: {}",
                    fps, tickrate, bunny_count
                );
                if let Ok(mut metrics) = GLOBAL_METRICS.write() {
                    metrics.fps = fps;
                    metrics.tps = tickrate;
                    metrics.bunnies = bunny_count;
                }
            }
        }
    }
}

#[wasm_bindgen]
pub fn get_rust_metrics() -> Result<JsValue, JsValue> {
    let metrics = GLOBAL_METRICS
        .read()
        .map_err(|_| JsValue::from_str("Failed to read metrics"))?;
    Ok(serde_wasm_bindgen::to_value(&*metrics)?)
}
