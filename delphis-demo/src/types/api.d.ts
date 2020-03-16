/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: GetDiscussionByID
// ====================================================

export interface GetDiscussionByID_discussion_moderator_userProfile_twitterURL {
  __typename: "URL";
  displayText: string;
  url: string;
}

export interface GetDiscussionByID_discussion_moderator_userProfile {
  __typename: "UserProfile";
  id: string;
  displayName: string;
  profileImageURL: string;
  twitterURL: GetDiscussionByID_discussion_moderator_userProfile_twitterURL;
}

export interface GetDiscussionByID_discussion_moderator {
  __typename: "Moderator";
  id: string;
  userProfile: GetDiscussionByID_discussion_moderator_userProfile;
}

export interface GetDiscussionByID_discussion_participants {
  __typename: "Participant";
  participantID: number | null;
}

export interface GetDiscussionByID_discussion_posts {
  __typename: "Post";
  id: string;
  isDeleted: boolean;
  deletedReasonCode: PostDeletedReason | null;
  content: string;
}

export interface GetDiscussionByID_discussion {
  __typename: "Discussion";
  id: string;
  anonymityType: AnonymityType;
  moderator: GetDiscussionByID_discussion_moderator;
  participants: GetDiscussionByID_discussion_participants[] | null;
  posts: GetDiscussionByID_discussion_posts[] | null;
  title: string;
}

export interface GetDiscussionByID {
  discussion: GetDiscussionByID_discussion | null;
}

export interface GetDiscussionByIDVariables {
  discussionID: string;
}

/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: ListDiscussions
// ====================================================

export interface ListDiscussions_listDiscussions_moderator_userProfile_twitterURL {
  __typename: "URL";
  displayText: string;
  url: string;
}

export interface ListDiscussions_listDiscussions_moderator_userProfile {
  __typename: "UserProfile";
  id: string;
  displayName: string;
  profileImageURL: string;
  twitterURL: ListDiscussions_listDiscussions_moderator_userProfile_twitterURL;
}

export interface ListDiscussions_listDiscussions_moderator {
  __typename: "Moderator";
  id: string;
  userProfile: ListDiscussions_listDiscussions_moderator_userProfile;
}

export interface ListDiscussions_listDiscussions_participants {
  __typename: "Participant";
  participantID: number | null;
}

export interface ListDiscussions_listDiscussions_posts {
  __typename: "Post";
  id: string;
  isDeleted: boolean;
  deletedReasonCode: PostDeletedReason | null;
  content: string;
}

export interface ListDiscussions_listDiscussions {
  __typename: "Discussion";
  id: string;
  anonymityType: AnonymityType;
  moderator: ListDiscussions_listDiscussions_moderator;
  participants: ListDiscussions_listDiscussions_participants[] | null;
  posts: ListDiscussions_listDiscussions_posts[] | null;
  title: string;
}

export interface ListDiscussions {
  listDiscussions: ListDiscussions_listDiscussions[] | null;
}

/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: CreatePost
// ====================================================

export interface CreatePost_addPost_moderator_userProfile_twitterURL {
  __typename: "URL";
  displayText: string;
  url: string;
}

export interface CreatePost_addPost_moderator_userProfile {
  __typename: "UserProfile";
  id: string;
  displayName: string;
  profileImageURL: string;
  twitterURL: CreatePost_addPost_moderator_userProfile_twitterURL;
}

export interface CreatePost_addPost_moderator {
  __typename: "Moderator";
  id: string;
  userProfile: CreatePost_addPost_moderator_userProfile;
}

export interface CreatePost_addPost_participants {
  __typename: "Participant";
  participantID: number | null;
}

export interface CreatePost_addPost_posts {
  __typename: "Post";
  id: string;
  isDeleted: boolean;
  deletedReasonCode: PostDeletedReason | null;
  content: string;
}

export interface CreatePost_addPost {
  __typename: "Discussion";
  id: string;
  anonymityType: AnonymityType;
  moderator: CreatePost_addPost_moderator;
  participants: CreatePost_addPost_participants[] | null;
  posts: CreatePost_addPost_posts[] | null;
  title: string;
}

export interface CreatePost {
  addPost: CreatePost_addPost;
}

export interface CreatePostVariables {
  discussionID: string;
  postContent: string;
}

/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL fragment: discussionInfo
// ====================================================

export interface discussionInfo_moderator_userProfile_twitterURL {
  __typename: "URL";
  displayText: string;
  url: string;
}

export interface discussionInfo_moderator_userProfile {
  __typename: "UserProfile";
  id: string;
  displayName: string;
  profileImageURL: string;
  twitterURL: discussionInfo_moderator_userProfile_twitterURL;
}

export interface discussionInfo_moderator {
  __typename: "Moderator";
  id: string;
  userProfile: discussionInfo_moderator_userProfile;
}

export interface discussionInfo_participants {
  __typename: "Participant";
  participantID: number | null;
}

export interface discussionInfo_posts {
  __typename: "Post";
  id: string;
  isDeleted: boolean;
  deletedReasonCode: PostDeletedReason | null;
  content: string;
}

export interface discussionInfo {
  __typename: "Discussion";
  id: string;
  anonymityType: AnonymityType;
  moderator: discussionInfo_moderator;
  participants: discussionInfo_participants[] | null;
  posts: discussionInfo_posts[] | null;
  title: string;
}

/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL fragment: profileInfo
// ====================================================

export interface profileInfo_twitterURL {
  __typename: "URL";
  displayText: string;
  url: string;
}

export interface profileInfo {
  __typename: "UserProfile";
  id: string;
  displayName: string;
  twitterURL: profileInfo_twitterURL;
}

/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

export enum AnonymityType {
  STRONG = "STRONG",
  UNKNOWN = "UNKNOWN",
  WEAK = "WEAK",
}

export enum PostDeletedReason {
  MODERATOR_REMOVED = "MODERATOR_REMOVED",
  PARTICIPANT_REMOVED = "PARTICIPANT_REMOVED",
  UNKNOWN = "UNKNOWN",
}

//==============================================================
// END Enums and Input Objects
//==============================================================
