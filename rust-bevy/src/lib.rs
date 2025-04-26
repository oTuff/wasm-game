use bevy::{
    asset::AssetMetaCheck,
    diagnostic::{DiagnosticsStore, FrameTimeDiagnosticsPlugin},
    prelude::*,
    winit::{UpdateMode, WinitSettings},
};
use fastrand;
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
        .add_systems(Update, initialize_bunny_image_data)
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
    commands.spawn((Text::new("FPS:"), ColorText));
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
    }
    spawn_bunnies(10, commands, bunny_image.into());
}

fn spawn_bunnies(amount: u32, mut commands: Commands, bunny_image: Res<BunnyImage>) {
    // let bunny_handle = &bunny_image.handle;
    // let bunny_handle = bunny_image.handle.clone();

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
                // image: bunny_handle.clone(),
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

fn update_positions(
    // time: Res<Time>,
    mut bunny_query: Query<(&mut Velocity, &mut Transform), With<Bunny>>,
) {
    for (mut velocity, mut transform) in bunny_query.iter_mut() {
        transform.translation.x += velocity.x; // * time.delta_secs();
        transform.translation.y += velocity.y; // * time.delta_secs();
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

// TODO: implement
// fn input_system()

// Debug information
fn update_diagnostics_text(
    diagnostics: Res<DiagnosticsStore>,
    mut writer: TextUiWriter,
    text_query: Query<Entity, With<ColorText>>,
) {
    if let Some(fps_diag) = diagnostics.get(&FrameTimeDiagnosticsPlugin::FPS) {
        if let Some(fps) = fps_diag.smoothed() {
            for entity in &text_query {
                *writer.text(entity, 0) = format!("FPS: {:.0}", fps);
            }
        }
    }
}
