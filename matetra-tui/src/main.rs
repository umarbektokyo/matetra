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

    let inner_p1 = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::Green))
        .block(
            Block::default()
                .borders(Borders::ALL)
                .title("team info block")
                .border_type(BorderType::Rounded),
        );

    let inner_p2 = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::LightCyan))
        .block(
            Block::default()
                .borders(Borders::ALL)
                .title("card block")
                .border_type(BorderType::Rounded),
        );

    frame.render_widget(inner_p1, inner_layout[0]);
    frame.render_widget(inner_p2, inner_layout[1]);

    let card_type = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::LightCyan))
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded),
        );

    let core_thingy = Paragraph::new("")
        .style(Style::default().fg(ratatui::style::Color::LightCyan))
        .block(
            Block::default()
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded),
        );

    let a_f_layout = Layout::default()
        .constraints(vec![
            Constraint::Percentage(10),
            Constraint::Percentage(50),
            Constraint::Percentage(40),
        ])
        .split(inner_layout[1]);

    frame.render_widget(card_type, a_f_layout[0]);
    frame.render_widget(core_thingy, a_f_layout[2]);
}
