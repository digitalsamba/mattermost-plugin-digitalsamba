export type UserConfig = {
    naming_scheme: string;
    embedded: boolean;
    show_prejoin_page: boolean;
}

export type MeetingInfo = {
    meeting_id: string;
    room_id: string;
    room_url: string;
    token: string;
    team_name: string;
    room_name: string;
}

export type DigitalSambaState = {
    userConfig: UserConfig | null;
    embeddedMeetings: MeetingInfo[];
}

export type MeetingConfig = {
    roomId: string;
    roomUrl: string;
    token: string;
    embedded: boolean;
    showPrejoinPage: boolean;
}