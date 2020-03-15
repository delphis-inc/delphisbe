import { gql, DocumentNode } from "apollo-boost";
import DiscussionFragments from '../../fragments/discussion';

const discussionQuery: DocumentNode = gql`
    query GetDiscussionByID($discussionID: ID!) {
        discussion(id: $discussionID) {
            ...discussionInfo
        }
    }
    ${DiscussionFragments.discussionInfo}`

const discussionListQuery: DocumentNode = gql`
    query ListDiscussions {
        listDiscussions {
            ...discussionInfo
        }
    }
    ${DiscussionFragments.discussionInfo}`

export { discussionQuery, discussionListQuery }