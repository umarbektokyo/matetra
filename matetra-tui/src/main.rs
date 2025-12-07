use std::vec;

use color_eyre::Result;
use crossterm::event::{self, Event};
use ratatui::{
    DefaultTerminal, Frame,
    layout::{Constraint, Layout},
    style::Style,
    widgets::{Block, BorderType, Borders, Paragraph},
};

fn main() -> Result<()> {
    color_eyre::install()?;
    let terminal = ratatui::init();
    let result = run(terminal);
    ratatui::restore();
    result
}

fn run(mut terminal: DefaultTerminal) -> Result<()> {
    loop {
        terminal.draw(render)?;
        if matches!(event::read()?, Event::Key(_)) {
            break Ok(());
        }
    }
}

fn render(frame: &mut Frame) {
    let outer_block = Block::default()
        .borders(Borders::ALL)
        .title("Matetra")
        .border_type(BorderType::Rounded);

    let outer_area = frame.area();
    frame.render_widget(outer_block.clone(), outer_area);

    let inner_area = outer_block.inner(outer_area);

    let inner_layout = Layout::default()
        .direction(ratatui::layout::Direction::Horizontal)
        .constraints(vec![Constraint::Percentage(60), Constraint::Percentage(40)])
        .split(inner_area);

    let team_info_block = Block::new()
        .style(Style::default().fg(ratatui::style::Color::Green))
        .borders(Borders::ALL)
        .border_type(BorderType::Rounded);

    let card_info_block = Block::new()
        .style(Style::default().fg(ratatui::style::Color::LightMagenta))
        .borders(Borders::ALL)
        .border_type(BorderType::Rounded);

    frame.render_widget(team_info_block, inner_layout[0]);
    frame.render_widget(card_info_block, inner_layout[1]);

    let card_upper_info = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::LightCyan))
        .block(Block::default());

    let card_bottom_block = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::LightCyan))
        .block(Block::default().borders(Borders::ALL));

    let card_info_block_layout = Layout::default()
        .constraints(vec![
            Constraint::Percentage(7),
            Constraint::Percentage(53),
            Constraint::Percentage(40),
        ])
        .split(inner_layout[1]);

    frame.render_widget(card_upper_info, card_info_block_layout[0]);
    frame.render_widget(card_bottom_block, card_info_block_layout[2]);

    let card_info_block_layout = Layout::default()
        .direction(ratatui::layout::Direction::Horizontal)
        .constraints(vec![
            Constraint::Percentage(20),
            Constraint::Percentage(65),
            Constraint::Percentage(15),
        ])
        .split(card_info_block_layout[0]);

    let card_info_top_left_part = Block::new()
        .borders(Borders::ALL)
        .style(Style::default().fg(ratatui::style::Color::LightRed));

    let card_info_top_left_part_paragraph = Paragraph::new("Constant")
        .style(Style::default())
        .alignment(ratatui::layout::Alignment::Center);

    let card_info_top_right_part_paragraph = Paragraph::new("Core")
        .style(Style::default())
        .alignment(ratatui::layout::Alignment::Center);

    let card_info_top_right_part = Block::new()
        .borders(Borders::ALL)
        .style(Style::default().fg(ratatui::style::Color::LightGreen));

    frame.render_widget(
        card_info_top_left_part_paragraph
            .clone()
            .block(card_info_top_left_part),
        card_info_block_layout[0],
    );
    frame.render_widget(
        card_info_top_right_part_paragraph
            .clone()
            .block(card_info_top_right_part),
        card_info_block_layout[2],
    );
}
