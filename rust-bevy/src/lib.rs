use bevy::prelude::*;
// use bevy_embedded_assets::EmbeddedAssetPlugin;
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn run_bevy_app() {
    App::new()
        // Set window resolution
        .add_plugins(DefaultPlugins.set(WindowPlugin {
            primary_window: Some(Window {
                resolution: (640., 480.).into(),
                ..default()
            }),
            ..default()
        }))
        // .add_plugins(EmbeddedAssetPlugin::default())
        .add_systems(Startup, setup)
        .run();
}

// #[derive(Component)]
// struct ColorText;

fn setup(mut commands: Commands, asset_server: Res<AssetServer>) {
    // UI camera
    commands.spawn(Camera2d);

    // commands.spawn(Camera2dBundle::default());

    commands.spawn(Sprite {
        image: asset_server.load("bunny.png"),
        ..default()
    });
    // // Text with one section
    // commands.spawn((
    //     // Accepts a `String` or any type that converts into a `String`, such as `&str`
    //     Text::new("hello\nbevy!"),
    //     TextFont {
    //         font_size: 67.0,
    //         ..default()
    //     },
    //     // Set the justification of the Text
    //     TextLayout::new_with_justify(JustifyText::Center),
    //     // Set the style of the Node itself.
    //     Node {
    //         position_type: PositionType::Absolute,
    //         bottom: Val::Px(5.0),
    //         right: Val::Px(5.0),
    //         ..default()
    //     },
    //     ColorText,
    // ));
}
